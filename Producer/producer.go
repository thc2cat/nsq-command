package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	nsq "github.com/nsqio/go-nsq"
)

func main() { // Producer

	type Unit struct {
		Date   string `json:"Date"`
		Sender string `json:"Sender"`
		Msg    string `json:"Msg"`
	}

	Hostname, _ := os.Hostname() // Get Short hostname

	sender := flag.String("s", Hostname, "sender")
	msg := flag.String("m", "", "message")
	msgFromFile := flag.String("M", "", "read message from file")

	nsqServer := flag.String("S", getenv("NSQ_SERVER", "localhost:4150"), "nsq server")
	nsqTopic := flag.String("T", "default", "nsq topic")
	verbose := flag.Bool("v", false, "verbose")

	flag.Parse()

	if *verbose {
		fmt.Printf("NSQ Producer set to [%s](%s)\n",
			*nsqServer, *nsqTopic)
	}

	w, _ := nsq.NewProducer(*nsqServer, nsq.NewConfig())

	if !*verbose {
		w.SetLogger(nil, 0) // zero logs
	}

	if *msgFromFile != "" {
		content, err := ioutil.ReadFile(*msgFromFile)
		if err != nil {
			log.Fatal(err)
		}
		*msg = (string)(content)
	}

	unit := &Unit{Sender: *sender, Date: time.Now().Format("15:04:05"), Msg: *msg}

	j, err := json.Marshal(unit)
	if err != nil {
		log.Panic("Could not connect")
	}
	err = w.Publish(*nsqTopic, j)
	if err != nil {
		log.Panic("Could not connect")
	}

	if *verbose {
		fmt.Printf("NSQ Producer sent.JSON [%s]\n", j)
	}

	w.Stop()

}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
