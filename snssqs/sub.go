package main

import (
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
		subscriptionName = flag.String("s", "", "Subscription name")
	)
	flag.Parse()

	if *subscriptionName == "" {
		log.Fatal("missing subscription name; " +
			"use -s to specify subscription name")
	}

	// Connect to PubSub system
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials("test", "test", ""),
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String(endpoints.EuCentral1RegionID),
		Endpoint:         aws.String("http://localhost:4566"),
	}))
	sqsService := sqs.New(sess)

	// Get Queue URL from Queue Name
	queueURL, err := urlFromQueueName(sqsService, *subscriptionName)
	if err != nil {
		log.Fatalf("unable to find/determine the URL of SQS queue %q: %v", *subscriptionName, err)
	}

	// Subscriber
	for {
		out, err := sqsService.ReceiveMessage(&sqs.ReceiveMessageInput{
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
			log.Printf("Recv: %s - %s\n", string(*m.Body), *m.MessageId)

			// Delete message (to ack)
			_, err = sqsService.DeleteMessage(&sqs.DeleteMessageInput{
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
