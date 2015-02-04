package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	MongoURI   = "localhost:27017"
	Database   = "hmon"
	Collection = "scans"

	Port = ":8080"
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

func main() {
	session, err := mgo.Dial(MongoURI)

	if err != nil {
		panic(err)
	}

	domains := ReadLines("../../data/topsites.txt")
	headers := []string{"x-xss-protection", "content-security-policy", "x-frame-options", "strict-transport-security", "nosniff", "x-powered-by", "server", "x-permitted-cross-domain-policies"}
	for {
		s := GenerateRandomScan(domains, headers)
		fmt.Println(s)
		collection := session.DB(Database).C(Collection)
		collection.Insert(s)
	}

}

func GenerateRandomScan(domains, headers []string) *Scan {
	s := new(Scan)
	s.TS = RandomDate()
	s.Domain = PickRandom(domains)

	for i := 0; i < rand.Intn(5); i++ {
		s.Headers = append(s.Headers, PickRandom(headers))
		s.Values = append(s.Values, []string{"fake"})
	}

	return s
}

func RandomDate() int {
	big := int(time.Now().Unix())
	small := int(time.Now().AddDate(0, -1, 0).Unix())
	return rand.Intn(big-small) + small
}

func PickRandom(a []string) string {
	return a[rand.Intn(len(a))]
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
