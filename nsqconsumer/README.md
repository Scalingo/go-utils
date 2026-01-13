# Package `nsqconsumer` v1.5.3

## Telemetry

Telemetry is enabled by default, and can be disabled with `WithoutTelemetry`
set to true in initialization options.

It records the duration of each handling (in seconds). The metric uses the
`scalingo.nsq.topic`, `scalingo.nsq.channel`, `scalingo.nsq.message_type`, and
`scalingo.nsq.status` attributes (`success` or `error`).

Metrics:

- `scalingo.nsq_consumer.message.duration`: handling duration in seconds
