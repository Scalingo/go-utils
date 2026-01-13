# Package `nsqproducer` v1.3.1

`nsqproducer` is a private package used by `nsqlbproducer`. It should **NEVER** be directly used! Please use the `nsqlbproducer` package.

`nsqproducer` contains the code to publish a message to a single `nsqd` instance. This is not good for high availability.

## Telemetry

Telemetry is enabled by default, and can be disabled with `WithoutTelemetry`
set to true in initialization options.

It records the number of published messages, the number of publish errors, and
the duration of each publish (in seconds). All metrics use the
`scalingo.nsq.topic`, `scalingo.nsq.message_type`, and
`scalingo.nsq.publish_type` attributes. The publish type is `immediate` for
`Publish`, and `deferred` for `DeferredPublish`.

Metrics:

- `scalingo.nsq_producer.publish.count`: number of published messages
- `scalingo.nsq_producer.publish.errors`: number of publish errors
- `scalingo.nsq_producer.publish.duration`: publish duration in seconds
