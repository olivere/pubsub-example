package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var (
		topicName = flag.String("t", "", "Topic name")
	)
	flag.Parse()

	if *topicName == "" {
		log.Fatal("missing topic name; " +
			"use -t to specify topic name")
	}

	// Create an SNS client
	//
	// The client is short running, so there's no need for e.g. a shutdown
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials("test", "test", ""),
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String(endpoints.EuCentral1RegionID),
		Endpoint:         aws.String("http://localhost:4566"),
	}))
	snsClient := sns.New(sess)

	// Create the topic if it doesn't exist yet
	topicARN, err := getTopicARN(snsClient, *topicName)
	if err == ErrTopicNotFound {
		topicARN, err = createTopic(snsClient, *topicName)
		if err != nil {
			log.Fatalf("unable to create SNS topic %s: %v", *topicName, err)
		}
	} else if err != nil {
		log.Fatalf("unable to find SNS topic %s: %v", *topicName, err)
	}

	// Publish some messages, synchronously
	var i int64
	for {
		i++
		now := time.Now().Format(time.RFC3339)
		payload := fmt.Sprintf("%08d %s", i, now)
		fmt.Printf("Send: %s\r", payload)
		out, err := snsClient.Publish(&sns.PublishInput{
			TopicArn: &topicARN,
			Message:  aws.String(payload),
		})
		if err != nil {
			log.Fatalf("unable to send message %s: %v", payload, err)
		}
		fmt.Printf("Send: %s - [%s]\n", payload, *out.MessageId)
		// Sleep for a random time of up to 1s
		time.Sleep(time.Duration(rand.Int63n(1000)) * time.Millisecond)
	}
}
