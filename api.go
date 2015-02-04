package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

var HeaderShortMap = map[string]string{"xss": "x-xss-protection", "xfo": "x-frame-options", "csp": "content-security-policy"}

func ServeAPI() {
	r := mux.NewRouter()
	r.HandleFunc("/api/hosts/search/{search}", SearchHandler)
	r.HandleFunc("/api/scans/", ScansHandler)
	http.Handle("/api/", r)

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)

	http.ListenAndServe(Port, nil)
}

// Builds the ts query for duraction
// Duraction can be today, week, month, total.
// Default to today
func BuildTimeQuery(duration string) bson.M {
	now := time.Now()
	if duration == "week" {
		weekAgo := now.AddDate(0, 0, -7)
		return bson.M{
			"ts": bson.M{
				"$gt": weekAgo.Unix(),
			},
		}
	} else if duration == "month" {
		monthAgo := now.AddDate(0, 0, -7)
		return bson.M{
			"ts": bson.M{
				"$gt": monthAgo.Unix(),
			},
		}
	} else if duration == "total" {
		return bson.M{}
	} else { // TODAY
		yesterday := now.AddDate(0, 0, -1)
		return bson.M{
			"ts": bson.M{
				"$gt": yesterday.Unix(),
			},
		}
	}

	return bson.M{}
}

func BuildDomainQuery(domain string) bson.M {
	if domain != "" {
		return bson.M{
			"domain": domain,
		}
	}
	return bson.M{}
}

/// XXX Content-security-policy, X-CSP, etc
func BuildHeaderQuery(header string) bson.M {
	if header != "" {

		fullheader, ok := HeaderShortMap[header]
		if !ok {
			fullheader = header
		}

		return bson.M{
			"headers": fullheader,
		}
	}

	return bson.M{}
}

func MergeQueries(q ...bson.M) bson.M {
	out := bson.M{}

	for _, queries := range q {
		for k, v := range queries {
			out[k] = v
		}
	}
	return out
}

func SearchHandler(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	searchTerm := vars["search"]

	results := []Scan{}
	c := Session.DB(Database).C(Collection)
	fmt.Println("searching for ", searchTerm)
	err := c.Find(bson.M{"host": &bson.RegEx{Pattern: searchTerm, Options: "i"}}).Sort("-ts").Limit(30).All(&results)

	out := []Scan{}
	seen := make(map[string]bool)
	for _, c := range results {
		if seen[c.Host] {
			continue
		}

		seen[c.Host] = true
		out = append(out, c)
	}

	if err != nil {
		panic(err)
	}

	res.Header().Set("Content-Type", "application/json")
	fmt.Fprint(res, ToJSON(out))
}

func ScansHandler(res http.ResponseWriter, req *http.Request) {
	headerQuery := BuildHeaderQuery(req.URL.Query().Get("header"))
	domainQuery := BuildDomainQuery(req.URL.Query().Get("domain"))
	timeQuery := BuildTimeQuery(req.URL.Query().Get("duration"))
	query := MergeQueries(headerQuery, domainQuery, timeQuery)

	fmt.Println(query)

	results := []Scan{}
	err := Session.DB(Database).C(Collection).Find(query).All(&results)

	if err != nil {
		panic(err)
	}

	res.Header().Set("Content-Type", "application/json")
	fmt.Fprint(res, ToJSON(results))
}
