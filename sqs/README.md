# AWS SQS

AWS SQS illustrates a one-to-one scenario:

- Publisher publishes a message
- One of the subscribers gets a copy of the message

## Prerequisites

1. Localstack must be running (see `docker-compose.yml` file in parent directory)

## Publisher

Publish messages on the `messages` topic:

```sh
./pub -t messages
```

Notice that publisher will try to create the topic if it doesn't already exist.

## Consumers

Consumer that subscribes to publisher on topic `messages` as such:

```sh
./sub -t messages
```

Notice that consumer will try to create the topic if it doesn't already exist.

## Remarks

1. Notice that using SNS alone will only give us a one-to-one scenario:
   One message gets relayed to exactly one of the subscribers, i.e. not
   every subscriber gets a copy. You need to use SNS in combination with
   SQS to get both (see `snssqs` parent directory).
