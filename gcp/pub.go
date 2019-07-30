package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/pubsub"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var (
		ctx       = context.Background()
		projectID = flag.String("p", "", "Project ID")
		topicName = flag.String("t", "", "Topic name")
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
		log.Fatal("missing topic name; use -t to specify a topic name")
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

	// Publish some messages, synchronously
	var i int64
	for {
		i++
		now := time.Now().Format(time.RFC3339)
		payload := fmt.Sprintf("%08d %s", i, now)
		log.Printf("Send: %s\n", payload)
		result := topic.Publish(ctx, &pubsub.Message{
			// Payload
			Data: []byte(payload),
			// Optional attributes/metadata
			Attributes: map[string]string{
				"origin": "native",
			},
		})
		// Block until the result is returned, which will include
		// a server-generated ID for the published message. In a
		// production environment, this should be put in a goroutine
		// if necessary at all.
		id, err := result.Get(ctx)
		if err != nil {
			log.Fatalf("unable to send message %s: %v", payload, err)
		}
		_ = id
		// Sleep for a random time of up to 1s
		time.Sleep(time.Duration(rand.Int63n(1000)) * time.Millisecond)
	}
}
