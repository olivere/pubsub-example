package main

import (
	"encoding/base64"
	"flag"
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
		log.Fatal("missing topic name; " +
			"use -t to specify topic name")
	}

	// Connect to PubSub system
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

	// Subscriber
	for {
		out, err := client.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            &queueURL,
			MaxNumberOfMessages: aws.Int64(1),
			VisibilityTimeout:   aws.Int64(3),
			WaitTimeSeconds:     aws.Int64(10),
			AttributeNames: aws.StringSlice([]string{
				"@ts",
			}),
			MessageAttributeNames: aws.StringSlice([]string{
				"All",
			}),
		})
		if err != nil {
			log.Printf("error on recv: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}
		if len(out.Messages) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}
		for _, m := range out.Messages {
			data, err := base64.StdEncoding.DecodeString(*m.Body)
			if err != nil {
				log.Printf("error decoding payload: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}
			log.Printf("Recv: %s - %s\n", string(data), *m.MessageId)

			// Delete message (to ack)
			_, err = client.DeleteMessage(&sqs.DeleteMessageInput{
				QueueUrl:      &queueURL,
				ReceiptHandle: m.ReceiptHandle,
			})
			if err != nil {
				log.Printf("error deleting message %s: %v", *m.MessageId, err)
				time.Sleep(1 * time.Second)
				continue
			}
		}
	}
}
