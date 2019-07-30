# NATS Streaming Server

## Prerequisites

```
$ export NATS_SERVER_URL=localhost:4222
```

## Publisher

Publish messages on the `messages` topic:

```
$ ./pub -t messages
```

## Consumers

Consumer with a random subscription name:

```
$ ./sub -t messages
```

Consumer with a fixed subscription name:

```
$ ./sub -t messages -s subscriber-1
```
