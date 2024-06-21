package graceful

import (
	"log"
	"os"
	"os/exec"
	"testing"
)

func TestMain(m *testing.M) {
	buildServer()
	m.Run()
}

// buildServer builds the test server binary
// In this function we build and run a binary, rather than using "go run"
// When running a Go program with go run, the interrupt signal (SIGINT) is not propagated to the
// child process by default.
// If "go run" is used tests that rely on the child process receiving the signal (those where the
// SIGHUP signal has been sent will leave the child process running indefinitely.)
func buildServer() {
	err := exec.Command("go", "build", "-o", "./testdata/server", "./testdata/cmd/server").Run()
	if err != nil {
		log.Println("Fail to build test server", err)
		os.Exit(1)
	}
}
