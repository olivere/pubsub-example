# NATS Streaming Server

## Prerequisites

```sh
export NATS_SERVER_URL=localhost:4222
```

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
