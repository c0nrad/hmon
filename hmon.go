package main

import (
	"fmt"
	"time"

	"github.com/robfig/cron"
	"gopkg.in/mgo.v2"
)

const (
	MongoURI   = "localhost:27017"
	Database   = "hmon"
	Collection = "scans"

	Port = ":8080"
)

var Session *mgo.Session

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
