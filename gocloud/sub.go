package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"time"

	stdgcppubsub "cloud.google.com/go/pubsub"
	"gocloud.dev/pubsub"
	_ "gocloud.dev/pubsub/gcppubsub"
	_ "gocloud.dev/pubsub/kafkapubsub"
	_ "gocloud.dev/pubsub/mempubsub"
	_ "gocloud.dev/pubsub/natspubsub"
	_ "gocloud.dev/pubsub/rabbitpubsub"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var (
		ctx = context.Background()
	)
	flag.Parse()

	if flag.NArg() != 1 {
		log.Fatal("missing connection string; see https://gocloud.dev/howto/pubsub/subscribe/ for details")
	}
	connection := flag.Arg(0)

	// Create the subscription if it doesn't exist yet
	sub, err := pubsub.OpenSubscription(ctx, connection)
	if err != nil {
		log.Fatalf("unable to open subscription %q: %v", connection, err)
	}
	defer sub.Shutdown(ctx)

	// Use As to check for specific PubSub implementations
	var gsub *stdgcppubsub.Subscription
	if sub.As(&gsub) {
		// We found that sub is a Google PubSub Subscription
		gsub.ReceiveSettings.MaxOutstandingMessages = 10
	}

	// Consume messages by pulling them in via a closure/callback handler
	for {
		msg, err := sub.Receive(ctx)
		if err != nil {
			log.Fatalf("unable to receive message: %v", err)
		}
		msg.Ack()
		log.Printf("Recv: %s\n", string(msg.Body))
	}
}
