# AWS SNS/SQS

AWS SNS/SQS illustrates a fan-out (one-to-many) scenario:

- Publisher publishes a message
- Every subscriber gets a copy of the message

## Prerequisites

1. Localstack must be running (see `docker-compose.yml` file in parent directory)
2. SQS Queues must be set up
3. SNS Topics must be set up and properly connected to SQS Queues
   (see [here](https://docs.aws.amazon.com/sns/latest/dg/sns-sqs-as-subscriber.html)
   and [here](https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-configure-subscribe-queue-sns-topic.html)).

Here's a short version of configuring Localstack with an SNS Topic `messages`
(for the publisher) and two SQS queues that will subscribe to the SNS Topic
`messages` (for the subscribers).

```sh
$ awslocal sns create-topic --name=messages
{
    "TopicArn": "arn:aws:sns:us-east-1:000000000000:messages"
}

$ awslocal sqs create-queue --queue-name=subscriber-1
{
    "QueueUrl": "http://localhost:4566/000000000000/subscriber-1"
}

$ awslocal sns subscribe \
   --topic-arn=arn:aws:sns:us-east-1:000000000000:messages \
   --protocol=sqs \
   --notification-endpoint=http://localhost:4566/000000000000/subscriber-1
{
    "SubscriptionArn": "arn:aws:sns:us-east-1:000000000000:messages:638eb44d-4b4b-4f02-93ac-b5ae097aa346"
}
```

Now the same with a 2nd subscriber `subscriber-2`:

```sh
$ awslocal sqs create-queue --queue-name=subscriber-2
{
    "QueueUrl": "http://localhost:4566/000000000000/subscriber-2"
}

$ awslocal sns subscribe \
   --topic-arn=arn:aws:sns:us-east-1:000000000000:messages \
   --protocol=sqs \
   --notification-endpoint=http://localhost:4566/000000000000/subscriber-2
{
    "SubscriptionArn": "arn:aws:sns:us-east-1:000000000000:messages:90cf06d6-5248-436a-bbc2-3a39000ca0e4"
}
```

To see that everyhing is properly set up, list the subscriptions to topic
`messages`:

```sh
$ awslocal sns list-subscriptions-by-topic --topic-arn=arn:aws:sns:us-east-1:000000000000:messages
{
    "Subscriptions": [
        {
            "SubscriptionArn": "arn:aws:sns:us-east-1:000000000000:messages:638eb44d-4b4b-4f02-93ac-b5ae097aa346",
            "Owner": "",
            "Protocol": "sqs",
            "Endpoint": "http://localhost:4566/000000000000/subscriber-1",
            "TopicArn": "arn:aws:sns:us-east-1:000000000000:messages"
        },
        {
            "SubscriptionArn": "arn:aws:sns:us-east-1:000000000000:messages:90cf06d6-5248-436a-bbc2-3a39000ca0e4",
            "Owner": "",
            "Protocol": "sqs",
            "Endpoint": "http://localhost:4566/000000000000/subscriber-2",
            "TopicArn": "arn:aws:sns:us-east-1:000000000000:messages"
        }
    ]
}
```

## Publisher

Publish messages on the `messages` topic:

```sh
./pub -t messages
```

## Consumers

Start two consumers the specific subscription/subscription name:

```sh
./sub -t messages -s subscriber-1
./sub -t messages -s subscriber-2
```

## Cleanup

For completion, here's how to clean up:

```sh
awslocal sqs delete-queue --queue-url=http://localhost:4566/000000000000/subscriber-2
awslocal sqs delete-queue --queue-url=http://localhost:4566/000000000000/subscriber-1
awslocal sns delete-topic --topic-arn=arn:aws:sns:us-east-1:000000000000:messages
$ awslocal sns list-topics
{
    "Topics": []
}
$ awslocal sqs list-queues
```
