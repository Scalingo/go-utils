# Package `nsqproducer` v1.4.0

`nsqproducer` is a private package used by `nsqlbproducer`. It should **NEVER** be directly used! Please use the `nsqlbproducer` package.

`nsqproducer` contains the code to publish a message to a single `nsqd` instance. This is not good for high availability.

## Telemetry

Telemetry is enabled by default, and can be disabled with `WithoutTelemetry`
set to true in initialization options.

It records the duration of each publish (in seconds). All metrics use the
`scalingo.nsq.topic`, `scalingo.nsq.message_type`,
`scalingo.nsq.publish_type`, and `scalingo.nsq.status` attributes. The publish
type is `immediate` for `Publish`, and `deferred` for `DeferredPublish`.
`scalingo.nsq.status` is `success` or `error`.

Metrics:

- `scalingo.nsq_producer.publish.duration`: publish duration in seconds
