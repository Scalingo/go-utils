package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Scalingo/go-utils/errors/v3"
	"github.com/Scalingo/go-utils/graceful"
)

func main() {
	numServers := 1

	// default options
	options := []graceful.Option{
		graceful.WithWaitDuration(time.Minute),
		graceful.WithPIDFile("./testdata/server.pid"),
	}

	// customise options
	for _, arg := range os.Args[1:] {
		// Split the option at the first '='
		idx := strings.Index(arg, "=")
		if idx == -1 {
			continue
		}
		opt := arg[:idx]
		val := arg[idx+1:]

		switch opt {
		case "pid-file":
			options = append(options, graceful.WithPIDFile(val))
		case "wait-duration":
			timeoutI, _ := strconv.Atoi(val)
			options = append(options, graceful.WithWaitDuration(time.Duration(timeoutI)*time.Millisecond))
		case "num-servers":
			numServers, _ = strconv.Atoi(val)
			options = append(options, graceful.WithNumServers(numServers))
		}
	}

	ctx := context.Background()
	s := graceful.NewService(
		options...,
	)

	errChan := make(chan error, numServers)
	var wg sync.WaitGroup

	for i := 0; i < numServers; i++ {
		wg.Add(1)
		port := 9000 + i
		endpoint := "/"
		if i > 0 {
			endpoint = fmt.Sprintf("/%d", i)
		}

		router := http.NewServeMux()
		router.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
			sleepStr := r.URL.Query().Get("sleep")
			sleep, _ := strconv.Atoi(sleepStr)
			if sleep != 0 {
				time.Sleep(time.Duration(sleep) * time.Millisecond)
			}
		})

		go func(i int) {
			defer wg.Done()
			addr := fmt.Sprintf(":%d", port)
			log.Printf("Serving on :%s\n", addr)
			err := s.ListenAndServe(ctx, "tcp", addr, router)
			if err != nil {
				log.Println("I'm dead because of", err)
				errChan <- errors.Wrapf(ctx, err, "I'm dead because of")
			}
		}(i)
	}

	wg.Wait()
	close(errChan)

	var shutdownErr error
	for err := range errChan {
		if shutdownErr == nil {
			shutdownErr = err
		} else {
			shutdownErr = errors.Wrap(ctx, shutdownErr, err.Error())
		}
	}

	if shutdownErr != nil {
		log.Fatal(ctx, shutdownErr)
	}

}
