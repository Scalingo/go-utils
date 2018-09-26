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

func buildServer() {
	err := exec.Command("go", "build", "-i", "-o", "./test-fixtures/server", "./test-fixtures/cmd/server").Run()
	if err != nil {
		log.Println("Fail to build test server", err)
		os.Exit(1)
	}
}
