package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"graceful"
)

func simpleHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello world\n")

	// Long-running go routine
	graceful.Go(func() {
		deadline := time.NewTimer(time.Minute)
		for {
			select {
			case <-deadline.C:
				return
			default:
			}
		}
	})
}

func main() {
	http.HandleFunc("/hello", simpleHandler)
	svr := &http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(simpleHandler),
	}
	go func() {
		err := svr.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()

	// Example go routine tracking.
	graceful.Go(func() {
		deadline := time.NewTimer(time.Second * 5)
		for {
			select {
			case <-deadline.C:
				return
			default:
			}
		}
	})

	// Specify graceful order.
	graceful.In(httpShutdown(svr), graceful.GoRoutineTerminator())

	// Wait for signals from os.
	graceful.Wait()
}

func httpShutdown(svr *http.Server) graceful.TerminatorFunc {
	return func() {
		if err := svr.Shutdown(context.Background()); err != nil {
			fmt.Println(err)
		}
	}
}
