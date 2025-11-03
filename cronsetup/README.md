# Package `cronsetup` v1.4.0

The `cronsetup` package eases the configuration of distributed cron jobs backed by etcd.

## Goal

This package aims at implementing a distributed and fault tolerant cron in order to:

* Run an identical process on several hosts
* Each of these process instantiate a cron with the same rules
* Ensure only one of these processes executes an iteration of a job

## Examples

Examples are available in the `examples` package. One can execute it with:

```sh
go run ./examples
```

It requires a etcd to be running at `127.0.0.1:2379`.
