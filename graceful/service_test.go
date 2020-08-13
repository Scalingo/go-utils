package graceful

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type cmdAndOutput struct {
	cmd    *exec.Cmd
	output *bytes.Buffer
}

func TestService_Shutdown_WithoutRequest(t *testing.T) {
	for _, s := range []os.Signal{syscall.SIGINT, syscall.SIGTERM} {
		t.Run("Signal "+s.String(), func(t *testing.T) {
			cmdAndOutput := startProcess(t)
			cmd := cmdAndOutput.cmd
			defer ensureProcessKilled(t, cmd)
			cmd.Process.Signal(s)
			isStoppedAfter(t, cmdAndOutput, 50*time.Millisecond)
		})
	}
}

func TestService_Shutdown_WithRequest(t *testing.T) {
	for _, s := range []os.Signal{syscall.SIGINT, syscall.SIGTERM} {
		t.Run("Signal "+s.String(), func(t *testing.T) {
			cmdAndOutput := startProcess(t)
			cmd := cmdAndOutput.cmd
			defer ensureProcessKilled(t, cmd)

			time.Sleep(10 * time.Millisecond)

			errs := make(chan error)
			go func() {
				_, err := http.Get("http://localhost:9000/?sleep=200")
				errs <- err
			}()
			time.Sleep(10 * time.Millisecond)

			cmd.Process.Signal(s)
			isRunningAfter(t, cmdAndOutput, 100*time.Millisecond)
			isStoppedAfter(t, cmdAndOutput, 300*time.Millisecond)
			require.NoError(t, <-errs)
		})
	}
}

func TestService_Shutdown_WithTimeout(t *testing.T) {
	for _, s := range []os.Signal{syscall.SIGINT, syscall.SIGTERM} {
		t.Run("Signal "+s.String(), func(t *testing.T) {
			cmdAndOutput := startProcess(t, "100")
			cmd := cmdAndOutput.cmd
			defer ensureProcessKilled(t, cmd)

			time.Sleep(10 * time.Millisecond)

			// Request will but cut
			errs := make(chan error)
			go func() {
				_, err := http.Get("http://localhost:9000/?sleep=1000")
				errs <- err
			}()
			time.Sleep(10 * time.Millisecond)

			cmd.Process.Signal(s)
			isRunningAfter(t, cmdAndOutput, 50*time.Millisecond)
			isStoppedAfter(t, cmdAndOutput, 150*time.Millisecond)
			require.Error(t, <-errs)
		})
	}
}

func TestService_Restart(t *testing.T) {
	cmdAndOutput := startProcess(t)
	cmd := cmdAndOutput.cmd
	defer ensureProcessKilled(t, cmd)

	time.Sleep(10 * time.Millisecond)

	errs := make(chan error, 100)
	go func() {
		defer close(errs)
		for i := 0; i < 100; i++ {
			_, err := http.Get("http://localhost:9000/?sleep=20")
			errs <- err
			time.Sleep(10 * time.Millisecond)
		}
	}()

	cmd.Process.Signal(syscall.SIGHUP)

	for err := range errs {
		require.NoError(t, err)
	}
}

func startProcess(t *testing.T, args ...string) cmdAndOutput {
	cmd := exec.Command("./test-fixtures/server", args...)
	b := new(bytes.Buffer)
	cmd.Stdout = b
	cmd.Stderr = b
	require.NoError(t, cmd.Start())
	return cmdAndOutput{cmd: cmd, output: b}
}

func isRunningAfter(t *testing.T, co cmdAndOutput, d time.Duration) {
	checkProcessAfter(t, co, d, true)
}

func isStoppedAfter(t *testing.T, co cmdAndOutput, d time.Duration) {
	checkProcessAfter(t, co, d, false)
}

func checkProcessAfter(t *testing.T, co cmdAndOutput, d time.Duration, shouldBeAlive bool) {
	cmd := co.cmd
	w := make(chan *os.ProcessState)
	go func() {
		cmd.Wait()
		w <- cmd.ProcessState
		close(w)
	}()
	timeout := time.NewTimer(d)
	defer timeout.Stop()
	select {
	case <-timeout.C:
		if !shouldBeAlive {
			t.Errorf("process %v was up after %v, output: \n\n%v", cmd, d, co.output.String())
		}
	case st := <-w:
		if shouldBeAlive {
			t.Errorf("process %v is dead after %v, status: %v, output: \n\n%v", cmd.Args, d, st.Success(), co.output.String())
		}
	}
	<-w
}

func ensureProcessKilled(t *testing.T, cmd *exec.Cmd) {
	ensurePidFileProcessKilled(t)
	err := cmd.Process.Kill()
	if err != nil && !strings.Contains(err.Error(), "already finished") {
		require.NoError(t, err)
	}
}

func ensurePidFileProcessKilled(t *testing.T) {
	out, err := ioutil.ReadFile("./test-fixtures/server.pid")
	if err == nil {
		pid, err := strconv.Atoi(strings.TrimSpace(string(out)))
		require.NoError(t, err)
		process, err := os.FindProcess(pid)
		require.NoError(t, err)
		err = process.Kill()
		if err != nil && !strings.Contains(err.Error(), "already finished") {
			require.NoError(t, err)
		}
		err = os.Remove("./test-fixtures/server.pid")
		require.NoError(t, err)
	}
}
