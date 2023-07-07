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
	err := exec.Command("go", "build", "-o", "./testdata/server", "./testdata/cmd/server").Run()
	if err != nil {
		log.Println("Fail to build test server", err)
		os.Exit(1)
	}
}
