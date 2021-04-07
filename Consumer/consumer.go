package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	nsq "github.com/nsqio/go-nsq"
)

// TODO - Add the write file with "sendervalue" as filename to consumer

func main() { // Consumer

	type Unit struct { // Unit type data used in our messages
		Date   string `json:"Date"`
		Sender string `json:"Sender"`
		Msg    string `json:"Msg"`
	}

	var u Unit

	Hostname, _ := os.Hostname() // Get Short hostname

	// Set NSQ_SERVER with -S or NSQ_SERVER env
	nsqServer := flag.String("S", getenv("NSQ_SERVER", "localhost:4150"), "nsq server")
	nsqTopic := flag.String("T", "default", "nsq topic")
	nsqChannel := flag.String("C", Hostname, "nsq channel")

	cmd := flag.String("c", "", "command to exec")
	verbose := flag.Bool("v", false, "verbose")

	flag.Parse()

	if *verbose {
		fmt.Printf("NSQ Consumer config Server[%s](Topic:%s/Channel:%s)\n",
			*nsqServer, *nsqTopic, *nsqChannel)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	q, _ := nsq.NewConsumer(*nsqTopic, *nsqChannel, nsq.NewConfig())

	if !*verbose {
		q.SetLogger(nil, 0)
	}

	q.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		err := json.Unmarshal(message.Body, &u)
		if err != nil {
			fmt.Println("json.Unmarshal error :", err)
		}

		nsqAction(*cmd, u.Date, u.Sender, u.Msg, *verbose)

		return nil
	}))
	err := q.ConnectToNSQD(*nsqServer)
	if err != nil {
		log.Panic("Could not connect", *nsqServer)
	}
	wg.Wait() // for ever
}

func nsqAction(action, date, sender, data string, verbose bool) {
	switch action {
	case "":
		fmt.Printf("NSQ stamp[%s], from[%s], data[%s]\n",
			date, sender, data)
	default:
		out := fmt.Sprintf(action, data)
		if verbose {
			fmt.Printf("Will execute ->%s<-\n", out)
		}
		tryexec(out)
	}
}

func tryexec(mycmd string) {
	waits := getenv("TIMEOUT", "90")
	wait, err := strconv.Atoi(waits)
	if err != nil {
		// handle error
		fmt.Println(err)
		os.Exit(1)
	}
	// Create a new context and add a timeout to it
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(wait)*time.Second)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	cmdargs := strings.Split(mycmd, " ")
	head := cmdargs[0]

	// Create the command with our context
	cmd := exec.CommandContext(ctx, head, cmdargs[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	// This time we can simply use Output() to get the result.
	out, err := cmd.Output()

	// We want to check the context error to see if the timeout was executed.
	// The error returned by cmd.Output() will be OS specific based on what
	// happens when a process is killed.
	if ctx.Err() == context.DeadlineExceeded {
		log.Println("Command timed out:", os.Args)
		return
	}

	// If there's no context error, we know the command completed (or errored).
	fmt.Print(string(out))
	if err != nil {
		log.Println("Non-zero exit code:", err)
	}
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
