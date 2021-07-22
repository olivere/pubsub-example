package main

import (
	"bytes"
	"context"
	"flag"
	"log"
	"math/rand"
	"os"
	"time"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/pubsub"
	"github.com/google/uuid"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var (
		ctx              = context.Background()
		projectID        = flag.String("p", "", "Project ID")
		topicName        = flag.String("t", "", "Topic name")
		subscriptionName = flag.String("s", "", "Subscription name")
	)
	flag.Parse()

	if *projectID == "" {
		// Maybe set from environment variable?
		*projectID = os.Getenv("PROJECT_ID")
		if *projectID == "" {
			// If we're running in Google Cloud, we can read the PROJECT_ID
			// from the infrastructure metadata. We ignore errors here because
			// an empty PROJECT_ID is allowed when using the PubSub emulator.
			*projectID, _ = metadata.ProjectID()
		}
	}
	if *topicName == "" {
		log.Fatal("missing -t option for topic name")
	}
	if *subscriptionName == "" {
		*subscriptionName = "sub-" + uuid.NewString()
	}

	// Connect to PubSub system
	client, err := pubsub.NewClient(ctx, *projectID)
	if err != nil {
		log.Fatalf("unable to create client: %v", err)
	}
	defer client.Close()

	// Create the topic if it doesn't exist yet
	topic := client.Topic(*topicName)
	found, err := topic.Exists(ctx)
	if err != nil {
		log.Fatalf("unable to check if topic %s exists: %v", *topicName, err)
	}
	if !found {
		topic, err = client.CreateTopic(ctx, *topicName)
		if err != nil {
			log.Fatalf("unable to create topic %s: %v", *topicName, err)
		}
	}

	// Create the subscription if it doesn't exist yet
	sub := client.Subscription(*subscriptionName)
	found, err = sub.Exists(ctx)
	if err != nil {
		log.Fatalf("unable to check if subscription %s exists: %v", *subscriptionName, err)
	}
	if !found {
		sub, err = client.CreateSubscription(ctx, *subscriptionName, pubsub.SubscriptionConfig{
			Topic:       topic,
			AckDeadline: 20 * time.Second,
		})
		if err != nil {
			log.Fatalf("unable to create subscription %s on topic %s: %v", *subscriptionName, *topicName, err)
		}
		// Fine-tune the subscription
		sub.ReceiveSettings.MaxOutstandingMessages = 10
	}

	// Consume messages by pulling them in via a closure/callback handler
	cctx, cancel := context.WithCancel(ctx)
	err = sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		msg.Ack()
		log.Printf("Recv: %s\n", string(msg.Data))

		// When QUIT gets sent, we stop
		if bytes.Equal(msg.Data, []byte("QUIT")) {
			cancel()
		}
	})
	if err != nil {
		log.Fatalf("Receive function failed with %v", err)
	}
}
