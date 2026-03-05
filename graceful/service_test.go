package graceful

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptrace"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// getCmd returns a command to run the server
func getCmd(args ...string) *exec.Cmd {
	return exec.Command("./testdata/server", args...)
}

type requestHandle struct {
	started <-chan struct{}
	errs    <-chan error
}

func startRequest(url string) requestHandle {
	started := make(chan struct{})
	errs := make(chan error, 2)

	go func() {
		defer close(errs)
		var startedOnce sync.Once
		closeStarted := func() {
			startedOnce.Do(func() {
				close(started)
			})
		}
		deadline := time.Now().Add(300 * time.Millisecond)

		for {
			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				errs <- err
				return
			}
			trace := &httptrace.ClientTrace{
				GotConn: func(httptrace.GotConnInfo) {
					closeStarted()
				},
			}
			req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				if time.Now().After(deadline) {
					errs <- err
					return
				}
				time.Sleep(5 * time.Millisecond)
				continue
			}

			errs <- nil
			errs <- resp.Body.Close()
			return
		}
	}()

	return requestHandle{
		started: started,
		errs:    errs,
	}
}

func waitForStart(t *testing.T, started <-chan struct{}) {
	t.Helper()

	const requestStartTimeout = 250 * time.Millisecond

	select {
	case <-started:
	case <-time.After(requestStartTimeout):
		t.Fatalf("request did not start after %v", requestStartTimeout)
	}
}

func requireNoRequestErrors(t *testing.T, errs <-chan error) {
	t.Helper()
	for err := range errs {
		require.NoError(t, err)
	}
}

func requireRequestError(t *testing.T, errs <-chan error) {
	t.Helper()
	hasErr := false
	for err := range errs {
		if err != nil {
			hasErr = true
		}
	}
	require.Truef(t, hasErr, "expected at least one request error")
}

// TestService_Shutdown_WithoutRequest tests the shutdown of the service without any request
func TestService_Shutdown_WithoutRequest(t *testing.T) {
	upgradeTimeout := time.Millisecond * 200
	shutdownTimeout := time.Millisecond * 100

	for i, s := range []os.Signal{syscall.SIGINT, syscall.SIGTERM} {
		t.Run("Send signal "+s.String()+" and expect service to stop", func(t *testing.T) {
			// Configure isGraceful
			isGraceful := newCmdAndOutput(t,
				withCmd(getCmd()),
				withUpgradeWaitDuration(upgradeTimeout),
				withShutdownWaitDuration(shutdownTimeout),
				withPidFile(fmt.Sprintf("./testdata/server-%d.pid", i)),
			)

			// start the command
			isGraceful.start()
			defer isGraceful.stop()

			// Send the signal
			isGraceful.signal(s)
			isGraceful.isStoppedAfter(shutdownTimeout)

			// Check the output
			output := isGraceful.getOutput()
			require.Containsf(t, output, "HTTP server is stopped", "OUTPUT:\n%v", output)
			require.Containsf(t, output, "No more connection running", "OUTPUT:\n%v", output)
		})
	}
}

// TestService_Shutdown_WithRequest tests the shutdown of the service with a request
func TestService_Shutdown_WithRequest(t *testing.T) {
	upgradeTimeout := time.Millisecond * 200
	shutdownTimeout := time.Millisecond * 100

	for i, s := range []os.Signal{syscall.SIGINT, syscall.SIGTERM} {
		t.Run("signal "+s.String()+" expect service to stop", func(t *testing.T) {
			// Configure isGraceful
			isGraceful := newCmdAndOutput(t,
				withCmd(getCmd()),
				withUpgradeWaitDuration(upgradeTimeout),
				withShutdownWaitDuration(shutdownTimeout),
				withPidFile(fmt.Sprintf("./testdata/server-%d.pid", i)),
			)

			// start the command
			isGraceful.start()
			defer isGraceful.stop()

			req := startRequest("http://localhost:9000/?sleep=200")
			waitForStart(t, req.started)

			// Send the signal
			isGraceful.signal(s)
			isGraceful.isRunningAfterAsync(100 * time.Millisecond)
			isGraceful.isStoppedAfterAsync(300 * time.Millisecond)

			requireNoRequestErrors(t, req.errs)

			// Check the output
			output := isGraceful.getOutput()
			require.Containsf(t, output, "HTTP server is stopped", "OUTPUT:\n%v", output)
			require.Containsf(t, output, "No more connection running", "OUTPUT:\n%v", output)
		})
	}
}

// TestService_Shutdown_MultipleServers_WithoutRequest tests the shutdown of the service with a request
func TestService_Shutdown_MultipleServers_WithoutRequest(t *testing.T) {
	upgradeTimeout := time.Millisecond * 200
	shutdownTimeout := time.Millisecond * 100

	for i, s := range []os.Signal{syscall.SIGINT, syscall.SIGTERM} {
		t.Run("signal "+s.String()+" expect service to stop", func(t *testing.T) {
			// Configure isGraceful
			isGraceful := newCmdAndOutput(t,
				withCmd(getCmd("num-servers=2")),
				withUpgradeWaitDuration(upgradeTimeout),
				withShutdownWaitDuration(shutdownTimeout),
				withPidFile(fmt.Sprintf("./testdata/server-%d.pid", i)),
			)

			// start the command
			isGraceful.start()
			defer isGraceful.stop()

			// Send the signal
			isGraceful.signal(s)
			isGraceful.isStoppedAfter(shutdownTimeout)

			// Check the output
			output := isGraceful.getOutput()
			require.Containsf(t, output, "HTTP server is stopped", "OUTPUT:\n%v", output)
			require.Containsf(t, output, "No more connection running", "OUTPUT:\n%v", output)
		})
	}
}

// TestService_Shutdown_MultipleServers_WithRequest tests the shutdown of the service with a request
func TestService_Shutdown_MultipleServers_WithRequest(t *testing.T) {
	upgradeTimeout := time.Millisecond * 200
	shutdownTimeout := time.Millisecond * 100

	for i, s := range []os.Signal{syscall.SIGINT, syscall.SIGTERM} {
		t.Run("signal "+s.String()+" expect service to stop", func(t *testing.T) {
			// Configure isGraceful
			isGraceful := newCmdAndOutput(t,
				withCmd(getCmd("num-servers=2")),
				withUpgradeWaitDuration(upgradeTimeout),
				withShutdownWaitDuration(shutdownTimeout),
				withPidFile(fmt.Sprintf("./testdata/server-%d.pid", i)),
			)

			// start the command
			isGraceful.start()
			defer isGraceful.stop()

			reqMain := startRequest("http://localhost:9000/?sleep=200")
			reqSecond := startRequest("http://localhost:9000/1?sleep=200")
			waitForStart(t, reqMain.started)
			waitForStart(t, reqSecond.started)

			// Send the signal
			isGraceful.signal(s)
			isGraceful.isRunningAfterAsync(100 * time.Millisecond)
			isGraceful.isStoppedAfterAsync(300 * time.Millisecond)

			requireNoRequestErrors(t, reqMain.errs)
			requireNoRequestErrors(t, reqSecond.errs)

			// Check the output
			output := isGraceful.getOutput()
			require.Containsf(t, output, "HTTP server is stopped", "OUTPUT:\n%v", output)
			require.Containsf(t, output, "No more connection running", "OUTPUT:\n%v", output)
		})
	}
}

// TestService_Shutdown_WithTimeout tests the shutdown of the service with a request that takes too long
func TestService_Shutdown_WithTimeout(t *testing.T) {
	for i, s := range []os.Signal{syscall.SIGINT, syscall.SIGTERM} {
		t.Run("signal "+s.String(), func(t *testing.T) {
			// Configure isGraceful
			isGraceful := newCmdAndOutput(t,
				withCmd(getCmd("wait-duration=100")),
				withUpgradeWaitDuration(200*time.Millisecond),
				withShutdownWaitDuration(100*time.Millisecond),
				withPidFile(fmt.Sprintf("./testdata/server-%d.pid", i)),
			)

			// start the command
			isGraceful.start()
			defer isGraceful.stop()

			req := startRequest("http://localhost:9000/?sleep=1000")
			waitForStart(t, req.started)

			// Send the signal
			isGraceful.signal(s)
			isGraceful.isRunningAfterAsync(50 * time.Millisecond)
			isGraceful.isStoppedAfterAsync(150 * time.Millisecond)

			// Check the output
			output := isGraceful.getOutput()
			assert.Containsf(t, output, "I'm dead because of shutdown service", "OUTPUT:\n%v", output)

			// The request should be unexpectedly terminated.
			requireRequestError(t, req.errs)
		})
	}
}

// TestService_Restart tests the restart of the service by sending a SIGHUP signal
// whilst the service receiving multiple requests
func TestService_Restart(t *testing.T) {
	// Configure isGraceful
	isGraceful := newCmdAndOutput(t,
		withCmd(getCmd()),
		withUpgradeWaitDuration(100*time.Millisecond),
		withShutdownWaitDuration(50*time.Millisecond),
		withPidFile("./testdata/server.pid"),
	)

	// start the command
	isGraceful.start()
	defer isGraceful.stop()

	errs := make(chan error, 256)
	stopRequests := make(chan struct{})
	firstRequestStarted := make(chan struct{})
	var requestWG sync.WaitGroup
	requestWG.Add(1)
	go func() {
		defer requestWG.Done()
		defer close(errs)
		started := false
		for {
			select {
			case <-stopRequests:
				return
			default:
			}

			if !started {
				close(firstRequestStarted)
				started = true
			}

			resp, err := http.Get("http://localhost:9000/?sleep=10")
			errs <- err
			if err == nil {
				// Response body must be closed
				err = resp.Body.Close()
				errs <- err
			}
		}
	}()

	waitForStart(t, firstRequestStarted)

	// Send the signal
	isGraceful.signal(syscall.SIGHUP)
	isGraceful.isRunningAfterAsync(30 * time.Millisecond)
	isGraceful.isRunningAfter(400 * time.Millisecond)

	close(stopRequests)
	requestWG.Wait()

	// The request should be no errors
	for err := range errs {
		require.NoError(t, err)
	}

	isGraceful.signal(syscall.SIGINT)
	isGraceful.isStoppedAfter(100 * time.Millisecond)

	// Check the output
	output := isGraceful.getOutput()
	require.Containsf(t, output, "Request graceful restart", "OUTPUT:\n%v", output)
	require.Containsf(t, output, "HTTP server is stopped", "OUTPUT:\n%v", output)
	require.Containsf(t, output, "No more connection running", "OUTPUT:\n%v", output)
}

type cmdAndOutput struct {
	t   *testing.T
	Cmd *exec.Cmd
	pid int

	waitGroup sync.WaitGroup

	output    *bytes.Buffer
	outputMu  sync.Mutex
	oldStdout io.Writer
	oldStderr io.Writer

	// shutdownWaitDuration is the duration which is waited for all connections to stop
	shutdownWaitDuration time.Duration

	// startWaitDuration is the duration to wait for a child process to start
	startWaitDuration time.Duration

	// upgradeWaitDuration is the duration the old process is waiting for
	// connection to close when a graceful restart has been ordered.
	upgradeWaitDuration time.Duration

	// pidFile tracks the pid of the last child among the chain of graceful restart
	pidFile string
}

// newCmdAndOutput creates a new cmdAndOutput struct using the functional options pattern
func newCmdAndOutput(t *testing.T, options ...func(*cmdAndOutput)) *cmdAndOutput {
	t.Helper()
	c := &cmdAndOutput{
		t:                    t,
		output:               new(bytes.Buffer),
		startWaitDuration:    500 * time.Millisecond,
		upgradeWaitDuration:  30 * time.Second,
		shutdownWaitDuration: 60 * time.Second,
	}
	for _, option := range options {
		option(c)
	}
	return c
}

// withCmd sets the Cmd field of the cmdAndOutput struct
func withCmd(cmd *exec.Cmd) func(*cmdAndOutput) {
	return func(c *cmdAndOutput) {
		c.Cmd = cmd
	}
}

// withPidFile sets the pidFile field of the cmdAndOutput struct
func withPidFile(pidFile string) func(*cmdAndOutput) {
	return func(c *cmdAndOutput) {
		c.pidFile = pidFile
	}
}

// withUpgradeWaitDuration sets the duration the old process is waiting for
func withUpgradeWaitDuration(duration time.Duration) func(output *cmdAndOutput) {
	return func(c *cmdAndOutput) {
		c.upgradeWaitDuration = duration
	}
}

// withShutdownWaitDuration sets the duration which is waited for all connections to stop
func withShutdownWaitDuration(duration time.Duration) func(output *cmdAndOutput) {
	return func(c *cmdAndOutput) {
		c.shutdownWaitDuration = duration
	}
}

// signal sends a signal to the process
func (c *cmdAndOutput) signal(signal os.Signal) {
	c.t.Helper()

	err := c.findProcess().Signal(signal)
	if err != nil {
		c.t.Fatalf("send signal %v: %v", signal, err)
	}
}

// start starts the process
func (c *cmdAndOutput) start() {
	c.t.Helper()

	c.oldStdout = c.Cmd.Stdout
	c.oldStderr = c.Cmd.Stderr
	r, w, _ := os.Pipe()
	c.Cmd.Stdout = w
	c.Cmd.Stderr = w

	// Read from pipe and append to buffer with locking
	go func() {
		b := make([]byte, 1024)
		for {
			n, err := r.Read(b)
			if n > 0 {
				c.outputMu.Lock()
				c.output.Write(b[:n])
				c.outputMu.Unlock()
			}
			if err != nil {
				break
			}
		}
	}()

	err := c.Cmd.Start()
	if err != nil {
		c.t.Fatalf("failed to start process: %v", err)
	}
	go func() {
		_ = c.Cmd.Wait()
	}()

	// Get the pid
	c.pid = c.Cmd.Process.Pid

	// Write the pid to the pid file
	if c.pidFile != "" {
		err := os.WriteFile(c.pidFile, []byte(strconv.Itoa(c.pid)), 0600)
		require.NoError(c.t, err)
	}

	c.waitForReady(c.startWaitDuration)
}

// stop stops the process
func (c *cmdAndOutput) stop() {
	// Wait for all (isRunningAfter / isStoppedAfter) operations to finish
	c.waitGroup.Wait()

	// send signal to parent process
	err := syscall.Kill(c.Cmd.Process.Pid, syscall.SIGTERM)
	if err != nil && !errors.Is(err, syscall.ESRCH) {
		c.t.Logf("kill process: %v", err)
	}

	// send signal to pid process
	err = syscall.Kill(c.currentPID(), syscall.SIGTERM)
	if err != nil && !errors.Is(err, syscall.ESRCH) {
		c.t.Logf("kill process: %v", err)
	}

	// Wait for the parent or child processes to finish
	c.isStoppedAfter(c.shutdownWaitDuration)

	// Delete pid file
	if c.pidFile != "" {
		err := os.Remove(c.pidFile)
		if err != nil && !os.IsNotExist(err) {
			require.NoError(c.t, err)
		}
	}
}

// isRunningAfter checks if the process is running after a certain duration
func (c *cmdAndOutput) isRunningAfter(timeout time.Duration) {
	c.t.Helper()
	c.checkProcessAfter(timeout, true)
}

// isRunningAfterAsync checks if the process is running after a certain duration, asynchronously
func (c *cmdAndOutput) isRunningAfterAsync(timeout time.Duration) {
	c.t.Helper()
	c.waitGroup.Add(1)
	go func() {
		defer c.waitGroup.Done()
		c.checkProcessAfter(timeout, true)
	}()
}

// isStoppedAfter checks if the process is stopped after a certain duration
func (c *cmdAndOutput) isStoppedAfter(timeout time.Duration) {
	c.t.Helper()
	c.checkProcessAfter(timeout, false)
}

// isStoppedAfterAsync checks if the process is stopped after a certain duration, asynchronously
func (c *cmdAndOutput) isStoppedAfterAsync(timeout time.Duration) {
	c.t.Helper()
	c.waitGroup.Add(1)
	go func() {
		defer c.waitGroup.Done()
		c.checkProcessAfter(timeout, false)
	}()
}

// checkProcessAfter checks the process is running after a certain duration
func (c *cmdAndOutput) checkProcessAfter(timeout time.Duration, shouldBeAlive bool) {
	c.t.Helper()

	// Has any process started
	require.NotNilf(c.t, c.Cmd.Process, "process %v hasn't started", c.Cmd)
	pid := c.currentPID()

	if shouldBeAlive {
		// Wait and then search for the process (parent or child)
		time.Sleep(timeout)
		p := c.findProcess()
		require.NoErrorf(c.t, p.Signal(syscall.Signal(0)), "process %v is dead after %v", pid, timeout)
	} else {
		// Poll for process disappearance to avoid waiting on non-child process.
		deadline := time.Now().Add(timeout)
		for {
			p := c.findProcess()
			err := p.Signal(syscall.Signal(0))
			if err != nil && (errors.Is(err, syscall.ESRCH) || errors.Is(err, os.ErrProcessDone)) {
				return
			}
			if time.Now().After(deadline) {
				c.t.Errorf("%v process %v was up after %v", time.Now(), pid, timeout)
				return
			}
			time.Sleep(5 * time.Millisecond)
		}
	}
}

func (c *cmdAndOutput) waitForReady(timeout time.Duration) {
	c.t.Helper()

	deadline := time.Now().Add(timeout)
	urls := c.readyURLs()
	for {
		allReady := true
		for _, u := range urls {
			resp, err := http.Get(u)
			if err != nil {
				allReady = false
				break
			}
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
		}
		if allReady {
			return
		}
		if time.Now().After(deadline) {
			c.t.Fatalf("server not ready after %v", timeout)
		}
		time.Sleep(5 * time.Millisecond)
	}
}

func (c *cmdAndOutput) readyURLs() []string {
	numServers := 1
	for _, arg := range c.Cmd.Args[1:] {
		if !strings.HasPrefix(arg, "num-servers=") {
			continue
		}
		value := strings.TrimPrefix(arg, "num-servers=")
		n, err := strconv.Atoi(value)
		if err == nil && n > 0 {
			numServers = n
		}
	}

	urls := make([]string, 0, numServers)
	for i := 0; i < numServers; i++ {
		if i == 0 {
			urls = append(urls, "http://localhost:9000/")
			continue
		}
		urls = append(urls, fmt.Sprintf("http://localhost:9000/%d", i))
	}

	return urls
}

// getOutput returns the output of the process
func (c *cmdAndOutput) getOutput() string {
	c.waitGroup.Wait()

	c.outputMu.Lock()
	defer c.outputMu.Unlock()
	return c.output.String()
}

func (c *cmdAndOutput) readPidFile() int {
	c.t.Helper()
	data, err := os.ReadFile(c.pidFile)
	require.NoError(c.t, err)
	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	require.NoError(c.t, err)
	return pid
}

func (c *cmdAndOutput) findProcess() *os.Process {
	pid := c.currentPID()
	p, err := os.FindProcess(pid)
	require.NoError(c.t, err)
	return p
}

func (c *cmdAndOutput) currentPID() int {
	if c.pidFile != "" {
		return c.readPidFile()
	}
	return c.pid
}
