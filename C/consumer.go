package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	nsq "github.com/nsqio/go-nsq"
)

func main() {

	type Unit struct {
		Date   string `json:"Date"`
		Sender string `json:"Sender"`
		Text   string `json:"Text"`
	}

	var unit Unit

	Hostname, _ := os.Hostname() // Get Short hostname

	nsqServer := flag.String("server", "193.51.24.101:4150", "nsq server")
	nsqTopic := flag.String("topic", "default", "nsq topic")
	nsqChannel := flag.String("channel", Hostname, "nsq channel")
	verbose := flag.Bool("v", false, "verbose")

	flag.Parse()

	if *verbose {
		fmt.Printf("NSQ Consumer set to [%s](%s/%s)\n",
			*nsqServer, *nsqTopic, *nsqChannel)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	nsqconfig := nsq.NewConfig()
	q, _ := nsq.NewConsumer(*nsqTopic, *nsqChannel, nsqconfig)

	if !*verbose {
		q.SetLogger(nil, 0)
	}

	q.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		err := json.Unmarshal(message.Body, &unit)
		if err != nil {
			fmt.Println(err)
		}

		nsqAction(unit.Date, unit.Sender, unit.Text)

		return nil
	}))
	err := q.ConnectToNSQD(*nsqServer)
	if err != nil {
		log.Panic("Could not connect", *nsqServer)
	}
	wg.Wait() // for ever
}

func nsqAction(date, sender, text string) {
	fmt.Printf("NSQ stamp[%s], from[%s], data[%s]\n",
		date, sender, text)
}
