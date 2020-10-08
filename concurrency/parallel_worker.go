package concurrency

import "sync"

// ParallelWorker is a struct providing an helper to run tasks in parallel
// while capping the parallelism with a known number of workers.
type ParallelWorker struct {
	sem         chan struct{}
	wg          *sync.WaitGroup
	endFunction func()
}

// NewParallelWorker constructs a new parallel worker running at maximum
// <workers> jobs at a time. <endFunc> is a callback called when all the jobs
// are over and that "CompleteProcessing" is called.
func NewParallelWorker(workers int, endFunc func()) ParallelWorker {
	w := ParallelWorker{
		sem:         make(chan struct{}, workers),
		wg:          &sync.WaitGroup{},
		endFunction: endFunc,
	}

	return w
}

// CompleteProcessing should be called to release all resources. It waits for all the runnings
// tasks to be over Calling Perform after this would result in a panic
func (w ParallelWorker) CompleteProcessing() {
	w.wg.Wait()
	close(w.sem)
	w.endFunction()
}

// Perform adds a job to handle. It can be called any number of time The
// concurrency will never get over the number of workers defined at the
// initialization
func (w ParallelWorker) Perform(function func()) {
	w.wg.Add(1)
	go func() {
		w.sem <- struct{}{}
		go func() {
			defer w.wg.Done()
			function()
			<-w.sem
		}()
	}()
}
