# Package `nsqconsumer` v1.5.3

## Telemetry

Telemetry is enabled by default, and can be disabled with `WithoutTelemetry`
set to true in initialization options.

It records the number of handled messages, the number of handled messages with
errors, and the duration of each handling (in seconds). All metrics use the
`scalingo.nsq.topic`, `scalingo.nsq.channel`, and `scalingo.nsq.message_type`
attributes.

Metrics:

- `scalingo.nsq_consumer.message.count`: number of handled messages
- `scalingo.nsq_consumer.message.errors`: number of handled messages with errors
- `scalingo.nsq_consumer.message.duration`: handling duration in seconds
