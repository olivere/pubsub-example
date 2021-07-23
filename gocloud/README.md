# Gocloud

Using PubSub via gocloud.dev.

## Prerequisites

To be able to use Google PubSub and NATS, use:

```sh
export PUBSUB_EMULATOR_HOST=localhost:8086
export PUBSUB_PROJECT_ID=demo
export NATS_SERVER_URL=localhost:4222
```

## Google PubSub

### Publisher

Publish messages on the `messages` topic:

```sh
./pub gcppubsub://demo/messages
```

### Consumers

Consumer with subscription name `subscription-1` o:

```sh
./sub gcppubsub://demo/messages/subscription-1
```

## NATS

### Publisher

Publish messages on the `myapp.messages` topic:

```sh
./pub nats://myapp.messages
```

### Consumers

Consumer as subject `myapp.*`:

```sh
./sub 'nats://myapp.*'
```

Read more about [NATS Subject-based messaging](https://nats-io.github.io/docs/developer/concepts/subjects.html).
