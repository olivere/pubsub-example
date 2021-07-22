package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var (
		topicName = flag.String("t", "", "Topic name")
	)
	flag.Parse()

	if *topicName == "" {
		log.Fatal("missing topic name; use -t to specify the topic name")
	}

	// Create an SQS client
	//
	// The client is short running, so there's no need for e.g. a shutdown
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials("test", "test", ""),
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String(endpoints.EuCentral1RegionID),
		Endpoint:         aws.String("http://localhost:4566"),
	}))
	client := sqs.New(sess)

	// Get Queue URL from Queue Name
	queueURL, err := urlFromQueueName(client, *topicName)
	if _, ok := IsAWSErrCode(err, sqs.ErrCodeQueueDoesNotExist); ok {
		// Create the queue
		queueURL, err = createQueue(client, *topicName)
		if err != nil {
			log.Fatalf("unable to create SQS queue %q: %v", *topicName, err)
		}
	} else if err != nil {
		// Some other kind of error
		log.Fatalf("unable to determine the URL of SQS queue %q: %v", *topicName, err)
	}

	// Publish some messages, synchronously
	var i int64
	for {
		i++
		now := time.Now().Format(time.RFC3339)
		payload := fmt.Sprintf("%08d %s", i, now)
		base64Payload := base64.StdEncoding.EncodeToString([]byte(payload))
		fmt.Printf("Send: %s\r", payload)
		in := &sqs.SendMessageInput{
			MessageBody: aws.String(base64Payload),
			QueueUrl:    &queueURL,
		}
		out, err := client.SendMessage(in)
		if err != nil {
			log.Fatalf("unable to send message %s: %v", payload, err)
		}
		fmt.Printf("Send: %s - [%s]\n", payload, *out.MessageId)
		// Sleep for a random time of up to 1s
		time.Sleep(time.Duration(rand.Int63n(1000)) * time.Millisecond)
	}
}
