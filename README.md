# PubSub Examples

This repositories illustrates a few things about PubSub systems, how
they work, how they differ, and how to implement them in Go.


## Prerequisites

Although also mentioned in the README files of the subfolders, you
should probably just start open at least 4 console windows and set
these environment variables for the rest of the examples.

These make [Google PubSub emulator](https://cloud.google.com/pubsub/docs/emulator)
work locally.

```
$ export PUBSUB_EMULATOR_HOST=localhost:8086
$ export PUBSUB_PROJECT_ID=demo
```

This one is for [NATS Streaming Server](https://github.com/nats-io/nats-streaming-server).

```
$ export NATS_SERVER_URL=localhost:4222
```

Also, make sure you have a recent version of Docker and Docker Compose
installed. We start the infrastructure through Docker Compose, so you
don't have to install anything besides Go on your machine.


## What is a PubSub system?

A PubSub system is a messaging system that has, as its name implies, two
components: Publisher of messages and subscriber to messages. In contrast
to synchronous communication, the publisher doesn't have to wait for a
message to be received, as well as the receiver doesn't have to be online
to retrieve messages sent earlier. As such, a PubSub system acts like a
buffer for asynchronous messaging.

Although there are various types of PubSub systems (Google PubSub, Kafka,
RabbitMQ, NATS Streaming Server to name just a few), they slightly differ in
how they work internally—and in terminology. So it is hard to find a
good abstraction for them, although libraries like
[`gocloud.dev/pubsub`](https://gocloud.dev/howto/pubsub/)
try to do just that. The `gocloud.dev/pubsub` is especially interesting as it
allows implementers to switch between different PubSub systems simply by using
a different URL. That is very welcome as you can e.g. use NATS in development
and Google PubSub in production, or a memory-based mock PubSub system 
while testing.

What all PubSub systems I've seen have in common is a way to publish messages
via a topic (although the naming can be different). That topic is simply a
name like `importer`, and a publisher send messages to that topic by some
client. The messages are received by the PubSub system and typically persisted
on the disk (some have in-memory storage as well). On the consumer side,
subscribers receive messages by creating a subscription on that topic, which
is also simply a name. Now that's where most of the abstractions start to fail.

Notice that some companies use PubSub systems like a storage solution. E.g.
the [New York Times uses Kafka for their publishing pipeline](https://www.confluent.io/blog/publishing-apache-kafka-new-york-times/).
They use Kafka as an append-only log where messages are never deleted.
You can iterate over the messages just like you can read through the rows
of a SQL database.

The most enlightening read about PubSub systems (or better: logs) is probably
"What every software engineer should know about real-time data's unifying abstraction"
by Jay Krebs from 2013, in which he lays out the fundamentals of Kafka:
https://engineering.linkedin.com/distributed-systems/log-what-every-software-engineer-should-know-about-real-time-datas-unifying
Highly recommended.


## Differences in PubSub systems

Some PubSub systems create topics on the fly, some systems require the caller
to create the topic before sending (or trying to receive) the first message.

On the subscriber side, some systems use a pull model (the consumer asks the
PubSub system) or a push model (the PubSub system forwards new messages), and
some offer both. On top of that, both models often have a synchronous or an
asynchronous implementation. Google PubSub, for example, has all of this: Push,
pull, sync/async.

What's also interesting is the guarantees that PubSub systems give you in terms
of ordering, delivery, and duplicates. In general, it is very hard to get
"exactly-once" delivery, and I think only Kafka can give this guarantee at the
time of writing this. So your code must be prepared for the case that you get
a message twice, or a message delivered later to come first to a subscriber.

Notice that some PubSub systems remove messages when they are consumed. Kafka,
on the other hand, doesn't do that. Instead, Kafka has a TTL for each message
and will remove them later. So setting the TTL to "never" will basically store
all messages forever, allowing you to iterate over historic messages years later.


## Publisher side

Why is PubSub used instead of synchronous messaging in e.g. a request/response style?

First of all PubSub is asynchronous by default. Suppose you want to track who
is visiting your web site on a per-page basis. Then, in your HTTP handler, you
simply gather the visitor details and send it to the PubSub system, in a kind
of fire-and-forget manner. You don't care if an event is lost, and you want your
HTTP handler to finish as quickly as possible. With request/response style, you
had to wait until the event system acknowledged your message.

Also, PubSub systems can send to more than one receiver with a single message.
We'll see later in "Subscriptions" how this can be used for different
architectural designs on the consumer side.


## Subscriber side

In PubSub systems, the publisher is more often the easier part of the equation.
What about the subscriber side?

When talking about PubSub systems with people new to the technology, there is
a hidden assumption that there is a 1:1 relationship between publisher and
subscriber. But that's not the case.

Let's get back to that example with the events sent by out HTTP handler.
Suppose it sends details about the endpoint, the visitor, and the time it took
to process the request. Wouldn't it be nice to have a loosely-coupled system
where one component tracks latency and throughput while a completely different
component tracks the pages that visitors are mostly interested in? With PubSub,
you can do that rather simply. E.g. with Google PubSub, if you subscribe with
to a topic with a random subscription name, you will get a copy of each and every
message sent to that topic. Repeat that: Every subscriber will get a copy.
This is called fan-out.

While the above is helpful for some work patterns, sometimes you don't want every
subscriber to get a copy of each message. Instead, you want to load-balance the
messages over a set of workers. Suppose you want to import products into a
data store, and want to be able to scale with the number of products being
imported concurrently. What you do with Google PubSub, for example, is to write
a publisher to sends products into an `importer` topic. On the listening side,
you create a set of subscribers that all have the same name, e.g. `importer_workers`.
What Google PubSub does when pushing/pulling a message to/from a subscriber,
it will only send one copy of the product to any of those `importer_workers`.
So if you want to scale with the load of products being imported, you simply start
more workers (a.k.a. subscribers). That is very easy to achieve, even automatically,
with e.g. a Kubernetes cluster.

Notice that most PubSub systems need the subscriber to acknowledge a message
that it successfully processed. If they forget to acknowledge, the message is
typically re-sent several times before giving up. Giving up, for e.g. RabbitMQ,
means to send those messages into a so-called dead letter queue.


## Examples

Now, this repository comes with 3 examples, each on in a subdirectory.

The [`gcp`](https://github.com/olivere/pubsub-example/tree/master/gcp)
and [`nats`](https://github.com/olivere/pubsub-example/tree/master/nats)
subdirectories illustrate how to implement a rather simply publisher and
subscriber for Google PubSub (via the Emulator) and NATS Streaming Server.

In [`gocloud`](https://github.com/olivere/pubsub-example/tree/master/gocloud)
you can see a single implementation that uses the wonderful
[`gocloud.dev` library by Google](https://github.com/google/go-cloud) to
do PubSub with a single implementation. You can use the code with any of the
supported PubSub systems: AWS SNS/SQS, Azure SB, Google PubSub, Kafka,
RabbitMQ, NATS Streaming Server, and an in-memory implementation which is
perfect for testing. You control the system to use by passing a URL.

Now, open 4 terminals and start the infrastructure in the 1st one:

```
$ docker-compose up
```

Head to the 2nd terminal and run:

```
$ cd gcp
$ make
$ ./pub -h
Usage of ./pub:
  -p string
    	Project ID
  -t string
    	Topic name
$ ./pub -t messages
2019/07/29 17:51:12 pub.go:66: Send: 00000001 2019-07-29T17:51:12+02:00
2019/07/29 17:51:12 pub.go:66: Send: 00000002 2019-07-29T17:51:12+02:00
2019/07/29 17:51:13 pub.go:66: Send: 00000003 2019-07-29T17:51:13+02:00
2019/07/29 17:51:14 pub.go:66: Send: 00000004 2019-07-29T17:51:14+02:00
2019/07/29 17:51:15 pub.go:66: Send: 00000005 2019-07-29T17:51:15+02:00
2019/07/29 17:51:15 pub.go:66: Send: 00000006 2019-07-29T17:51:15+02:00
...
```

Go to the 3rd terminal and run:

```
$ cd gcp
$ ./sub -h
Usage of ./sub:
  -p string
    	Project ID
  -s string
    	Subscription name
  -t string
    	Topic name
$ ./sub -t messages -s subscriber-1
2019/07/29 17:51:12 sub.go:87: Recv: 00000001 2019-07-29T17:51:12+02:00
2019/07/29 17:51:12 sub.go:87: Recv: 00000002 2019-07-29T17:51:12+02:00
2019/07/29 17:51:13 sub.go:87: Recv: 00000003 2019-07-29T17:51:13+02:00
2019/07/29 17:51:14 sub.go:87: Recv: 00000004 2019-07-29T17:51:14+02:00
2019/07/29 17:51:15 sub.go:87: Recv: 00000005 2019-07-29T17:51:15+02:00
2019/07/29 17:51:16 sub.go:87: Recv: 00000006 2019-07-29T17:51:15+02:00
```

Up to the 4th terminal and run:

```
$ cd gcp
$ ./sub -t messages -s subscriber-2
2019/07/29 17:51:12 sub.go:87: Recv: 00000001 2019-07-29T17:51:12+02:00
2019/07/29 17:51:12 sub.go:87: Recv: 00000002 2019-07-29T17:51:12+02:00
2019/07/29 17:51:13 sub.go:87: Recv: 00000003 2019-07-29T17:51:13+02:00
2019/07/29 17:51:14 sub.go:87: Recv: 00000004 2019-07-29T17:51:14+02:00
2019/07/29 17:51:15 sub.go:87: Recv: 00000005 2019-07-29T17:51:15+02:00
2019/07/29 17:51:16 sub.go:87: Recv: 00000006 2019-07-29T17:51:15+02:00
...
```

Notice how both `subscriber-1` and `subscriber-2` get a copy of each message
being sent by the publisher on topic `messages`.

Now, stop the `subscriber-2` in terminal 4, and run it under the same
name as `subscriber-1` instead:

```
$ ./sub -t messages -s subscriber-1
2019/07/29 17:52:05 sub.go:87: Recv: 00000007 2019-07-29T17:52:05+02:00
2019/07/29 17:52:07 sub.go:87: Recv: 00000010 2019-07-29T17:52:07+02:00
...
```

Notice how messages will be split between terminal 3 and 4, load-balanced if you will.

# License

MIT.
