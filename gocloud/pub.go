package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	rawgcp "cloud.google.com/go/pubsub/apiv1"
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

	// Connect to PubSub system
	topic, err := pubsub.OpenTopic(ctx, connection)
	if err != nil {
		log.Fatalf("unable to open topic %q: %v", connection, err)
	}
	defer topic.Shutdown(ctx)

	// Use As to check for specific PubSub implementations
	var gcppubclient *rawgcp.PublisherClient
	if topic.As(&gcppubclient) {
		// ...
	}

	// Publish some messages, synchronously
	var i int64
	for {
		i++
		now := time.Now().Format(time.RFC3339)
		payload := fmt.Sprintf("%08d %s", i, now)
		log.Printf("Send: %s\n", payload)
		err := topic.Send(ctx, &pubsub.Message{
			// Payload
			Body: []byte(payload),
			// Optional metadata
			Metadata: map[string]string{
				"origin": "native",
			},
		})
		if err != nil {
			log.Fatalf("unable to send message %s: %v", payload, err)
		}
		// Sleep for a random time of up to 1s
		time.Sleep(time.Duration(rand.Int63n(1000)) * time.Millisecond)
	}
}
