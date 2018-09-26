package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Scalingo/go-utils/graceful"
)

func main() {
	timeout := time.Minute
	if len(os.Args) == 2 && os.Args[1] != "" {
		timeoutI, _ := strconv.Atoi(os.Args[1])
		timeout = time.Duration(timeoutI) * time.Millisecond
	}
	ctx := context.Background()
	s := graceful.NewService(
		graceful.WithWaitDuration(timeout),
		graceful.WithPIDFile("./test-fixtures/server.pid"),
	)
	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		sleepStr := r.URL.Query().Get("sleep")
		sleep, _ := strconv.Atoi(sleepStr)
		if sleep != 0 {
			time.Sleep(time.Duration(sleep) * time.Millisecond)
		}
	})
	log.Println("Serving on :9000")
	err := s.ListenAndServe(ctx, "tcp", ":9000", router)
	if err != nil {
		log.Println("I'm dead because of", err)
		os.Exit(-1)
	}
}
