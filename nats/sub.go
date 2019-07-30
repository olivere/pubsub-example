package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/gofrs/uuid"
	stan "github.com/nats-io/stan.go"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var (
		url              = flag.String("url", stan.DefaultNatsURL, "NATS Server URLs, separated by commas")
		clusterID        = flag.String("cluster", "demo", "Cluster ID")
		clientID         = flag.String("client", "", "Client ID")
		topicName        = flag.String("t", "", "Topic name")
		subscriptionName = flag.String("s", "", "Subscription name")
	)
	flag.Parse()

	if *clientID == "" {
		*clientID = uuid.Must(uuid.NewV4()).String()
	}
	if *topicName == "" {
		log.Fatal("missing topic name; use -t to specify a topic name")
	}

	// Connect to NATS Streaming Server cluster
	sc, err := stan.Connect(*clusterID, *clientID,
		stan.NatsURL(*url),
		stan.Pings(10, 5),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Printf("Connection lost: %v", reason)
		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	var sub stan.Subscription
	if *subscriptionName != "" {
		// Subscribe to the topic as a queue.
		// Start with new messages as they come in; don't replay earlier messages.
		sub, err = sc.QueueSubscribe(*topicName, *subscriptionName, func(msg *stan.Msg) {
			msg.Ack()
			log.Printf("Recv: %s [topic=%s]\n", string(msg.Data), msg.Subject)
		}, stan.StartWithLastReceived())
		if err != nil {
			log.Fatal(err)
		}
	} else {
		sub, err = sc.Subscribe(*topicName, func(msg *stan.Msg) {
			msg.Ack()
			log.Printf("Recv: %s [topic=%s]\n", string(msg.Data), msg.Subject)
		}, stan.StartWithLastReceived())
		if err != nil {
			log.Fatal(err)
		}
	}

	// Wait for Ctrl+C
	doneCh := make(chan bool)
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt)
		<-sigCh
		sub.Unsubscribe()
		doneCh <- true
	}()
	<-doneCh
}
