package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"graceful"
)

func simpleHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello world\n")

	// Background go routine
	graceful.Go(func(ctx context.Context) {
		log.Println("Started long-running go routine from simpleHandler()")
		deadline := time.NewTimer(time.Second * 5)
		<-deadline.C
		to := time.Second * 5
		newCtx, cf := context.WithTimeout(ctx, to)
		graceful.Go(func(ctx context.Context) {
			defer cf()
			log.Println("Spawned inner go routine from simpleHandler()")
			k := 0
			for {
				select {
				case <-time.NewTicker(time.Second).C:
					log.Println("Ticking", k)
					k++
				case <-newCtx.Done():
					log.Println("Completed inner go routine from simpleHandler()", ctx.Err())
					return
				}
			}
		})
		log.Println("Completed long-running go routine from simpleHandler()", ctx.Err())
	})
}

func main() {
	http.HandleFunc("/hello", simpleHandler)
	svr := &http.Server{
		Addr:    ":8081",
		Handler: http.HandlerFunc(simpleHandler),
	}

	graceful.Go(
		func(ctx context.Context) {
			log.Println("Starting server")
			err := svr.ListenAndServe()
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				panic(err)
			}
		})
	graceful.Go(
		func(ctx context.Context) {
			<-ctx.Done()
			log.Println("Shutting down server")
			if err := svr.Shutdown(context.Background()); err != nil {
				log.Printf("Shutdown error: %v", err)
			}
		},
	)

	// Example go routine tracking.
	graceful.Go(func(ctx context.Context) {
		deadline := time.NewTimer(time.Second * 5)
		<-deadline.C
		log.Println("Completed running go routine from main().", ctx.Err())
	})

	// Wait for signals from os.
	graceful.Wait()
}
