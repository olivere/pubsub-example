package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/sqs"
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

func createQueue(service *sqs.SQS, queueName string) (string, error) {
	out, err := service.CreateQueue(&sqs.CreateQueueInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		return "", err
	}
	return *out.QueueUrl, nil
}

func IsAWSErrCode(err error, code string) (*awserr.Error, bool) {
	if aerr, ok := err.(awserr.Error); ok && aerr.Code() == code {
		return &aerr, true
	}
	return nil, false
}
