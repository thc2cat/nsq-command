package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	nsq "github.com/nsqio/go-nsq"
)

func main() { // Producer

	type Unit struct {
		Date   string `json:"Date"`
		Sender string `json:"Sender"`
		Text   string `json:"Text"`
	}

	Hostname, _ := os.Hostname() // Get Short hostname

	sender := flag.String("sender", Hostname, "sender")
	texte := flag.String("text", "", "text to send")

	nsqServer := flag.String("server", "193.51.24.101:4150", "nsq server")
	nsqTopic := flag.String("topic", "default", "nsq topic")
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

	unit := &Unit{Sender: *sender, Date: time.Now().Format("15:04:05"), Text: *texte}

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
