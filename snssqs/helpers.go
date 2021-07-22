package main

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var (
	ErrTopicNotFound = errors.New("topic not found")
)

func urlFromQueueName(service *sqs.SQS, queueName string) (string, error) {
	url, err := service.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		return "", err
	}
	return *url.QueueUrl, nil
}

func getTopicARN(service *sns.SNS, queueName string) (string, error) {
	in := &sns.ListTopicsInput{}
	for {
		out, err := service.ListTopics(in)
		if err != nil {
			return "", err
		}
		for _, topic := range out.Topics {
			if topic.String() == queueName {
				return *topic.TopicArn, nil
			}
		}
		if out.NextToken == nil {
			return "", ErrTopicNotFound
		}
		in.NextToken = out.NextToken
	}
}

func createTopic(service *sns.SNS, queueName string) (string, error) {
	out, err := service.CreateTopic(&sns.CreateTopicInput{
		Name: &queueName,
	})
	if err != nil {
		return "", err
	}
	return *out.TopicArn, nil
}

func IsAWSErrCode(err error, code string) (*awserr.Error, bool) {
	if aerr, ok := err.(awserr.Error); ok && aerr.Code() == code {
		return &aerr, true
	}
	return nil, false
}
