module github.com/olivere/google-pubsub-example

go 1.12

require (
	cloud.google.com/go v0.43.0
	github.com/gofrs/uuid v3.2.0+incompatible
	github.com/nats-io/nats-server/v2 v2.0.2 // indirect
	github.com/nats-io/nats-streaming-server v0.15.1 // indirect
	github.com/nats-io/stan.go v0.5.0
	github.com/prometheus/client_golang v0.9.3-0.20190127221311-3c4408c8b829 // indirect
	gocloud.dev v0.15.0
	gocloud.dev/pubsub/kafkapubsub v0.15.0
	gocloud.dev/pubsub/natspubsub v0.15.0
	gocloud.dev/pubsub/rabbitpubsub v0.15.0
)
