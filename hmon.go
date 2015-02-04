package main

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/robfig/cron"
	"gopkg.in/mgo.v2"
)

const (
	Collection = "scans"
)

var Session *mgo.Session
var Port = GetPort()
var Database = GetDBName()
var MongoURI = GetMongoURI()

func GetPort() string {
	port := os.Getenv("PORT")
	if port != "" {
		return ":" + port
	}
	return ":8080"
}

func GetDBName() string {
	mongoLabURI := os.Getenv("MONGOLAB_URI")
	if mongoLabURI == "" {
		return "hmon"
	}

	parsedUrl, _ := url.Parse(mongoLabURI)
	return parsedUrl.Path[1:]
}

func GetMongoURI() string {
	mongoLabURI := os.Getenv("MONGOLAB_URI")
	if mongoLabURI == "" {
		return "localhost:27017"
	}
	return mongoLabURI
}

func main() {
	session, err := mgo.Dial(MongoURI)

	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)
	Session = session

	go ScannerCron()
	ServeAPI()
}

func ScannerCron() {
	ScanHostFile("./data/http.top")

	c := cron.New()
	c.AddFunc("@hourly", func() {
		fmt.Println("Running scanner", time.Now())
		ScanHostFile("./data/http.top")
	})
	c.Start()
}
