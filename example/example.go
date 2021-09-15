package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	. "graceful"
)

func simpleHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello world\n")

	// Background go routine
	Go(func(ctx context.Context) {
		Go(func(ctx context.Context) {
			n := 0
			ticker := time.NewTicker(time.Second * 1)
			for {
				select {
				case <-ticker.C:
					n = n + 1
					log.Println("ticking", n)
				case <-ctx.Done():
					log.Println("completed ticking", ctx.Err())
					return
				}
			}
		})
	})
	Go(func(ctx context.Context) {
		log.Println("Running Long", "Started")
		<-time.NewTimer(time.Second * 10).C
		log.Println("Running Long", "Stopped", ctx.Err())
	})
}

func main() {
	svr := &http.Server{
		Addr:    ":8081",
		Handler: http.HandlerFunc(simpleHandler),
	}

	Go(func(ctx context.Context) {
		log.Println("Starting server")
		err := svr.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	})
	Go(func(ctx context.Context) {
		<-ctx.Done()
		log.Println("Shutting down server")
		if err := svr.Shutdown(context.Background()); err != nil {
			log.Printf("Shutdown error: %v", err)
		}
	})

	// Example go routine tracking.
	Go(func(ctx context.Context) {
		deadline := time.NewTimer(time.Second * 5)
		<-deadline.C
		log.Println("Completed running go routine from main().", ctx.Err())
	})

	// Wait for signals from os.
	Wait()
}
