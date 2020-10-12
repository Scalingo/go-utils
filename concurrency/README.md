# Package `concurrency`

## Parallel Worker

ParallelWorker is a struct providing an helper to run tasks in parallel while
capping the parallelism with a known number of workers.

```
endcallback := func() {}
w := NewParalleWorker(10, endcallback)

for _, item := range slice {
  w.Perform(func() {
    doSomething(item)
  })
}

// Wait for all jobs to be over
w.Stop()
```

Better example in `parallel_worker_example_test.go`
