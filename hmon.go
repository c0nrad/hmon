package main

import (
	"fmt"
	"os"
	"time"

	"github.com/robfig/cron"
	"gopkg.in/mgo.v2"
)

const (
	DefaultMongoURI = "localhost:27017"
	Database        = "hmon"
	Collection      = "scans"
)

var Session *mgo.Session
var Port = ":8080"

func main() {
	mongoURI := os.Getenv("MONGOLAB_URI")
	if mongoURI == "" {
		mongoURI = DefaultMongoURI
	}

	port := os.Getenv("PORT")
	if port != "" {
		Port = ":" + port
	}

	session, err := mgo.Dial(mongoURI)

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
