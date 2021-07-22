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

## Consumers

Consumer with a random subscription name:

```sh
./sub -t messages
```

Consumer with a fixed subscription name:

```sh
./sub -t messages -s subscriber-1
```

## Remarks

1. When using a unique subscription name, every subscriber gets a copy of the
   messages sent by the producer.
2. When sharing a subscription name, messages get sent to only one subscriber.
3. By default, subscribers receive only those messages sent after they've been
   started.
