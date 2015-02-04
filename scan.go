package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Scan struct {
	Id      bson.ObjectId `bson:"_id,omitempty"`
	Host    string        `bson:"host"`
	Domain  string        `bson:"domain"`
	TS      int           `bson:"ts"`
	Headers []string      `bson:"headers"`
	Values  [][]string    `bson:"values"`
	Error   string        `bson:"error"`
}

const (
	Threads = 1
)

func (s Scan) GetHeaders(search string) []string {
	for i, header := range s.Headers {
		if header == search {
			return s.Values[i]
		}
	}

	return []string{}
}

func (s Scan) Save() error {
	collection := Session.DB(Database).C(Collection)
	return collection.Insert(s)
}

func ToJSON(scans []Scan) string {
	b, err := json.Marshal(scans)
	if err != nil {
		return ""
	}
	return string(b)
}

func (s Scan) String() string {
	b, err := json.Marshal(s)
	if err != nil {
		return ""
	}
	return string(b)
}

func ScanHost(host string) *Scan {
	scan := new(Scan)
	scan.TS = int(time.Now().Unix())
	scan.Host = host
	scan.Domain = GrabDomain(host)

	req, err := http.Get(host)
	if err != nil {
		scan.Error = err.Error()
		return scan
	}

	for key, value := range req.Header {
		scan.Headers = append(scan.Headers, strings.ToLower(key))
		scan.Values = append(scan.Values, value)
	}

	return scan
}

func ScanHostFile(filename string) {
	hosts := ReadLines(filename)

	hostChan := make(chan string, Threads*10)
	wg := new(sync.WaitGroup)
	wg.Add(Threads)
	for i := 0; i < Threads; i++ {
		go ScanWorker(hostChan, wg)
	}

	for _, host := range hosts {
		hostChan <- host
	}
	wg.Wait()
}

func ScanWorker(hostChan chan string, wg *sync.WaitGroup) {
	for host := range hostChan {
		s := ScanHost(host)
		s.Save()
	}
	wg.Done()
}

func ReadLines(filename string) []string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(data), "\n")

	// Remove last line if empty
	if len(lines[len(lines)-1]) == 0 {
		lines = lines[0 : len(lines)-1]
	}
	return lines
}

func GrabDomain(u string) string {
	parsedUrl, _ := url.Parse(u)
	return parsedUrl.Host
}
