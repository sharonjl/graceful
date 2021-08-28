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
	graceful.Go(context.Background(), func(ctx context.Context) {
		log.Println("Started long-running go routine from simpleHandler()")
		deadline := time.NewTimer(time.Second * 5)
		<-deadline.C

		newCtx, cf := context.WithTimeout(ctx, time.Second*5)
		graceful.Go(newCtx, func(ctx context.Context) {
			defer cf()
			log.Println("Spawned inner go routine #1 from simpleHandler()")
			k := 0
			for {
				select {
				case <-time.NewTicker(time.Second).C:
					log.Println("Ticking", k)
					k++
				case <-ctx.Done():
					log.Println("Completed inner go routine #1 from simpleHandler()", ctx.Err())
					return
				}
			}
		})
		<-newCtx.Done()

		newCtx2, cf := context.WithCancel(ctx)
		graceful.Go(newCtx2, func(ctx context.Context) {
			log.Println("Spawned inner go routine #2 from simpleHandler()")
			k := 0
			for {
				select {
				case <-time.NewTicker(time.Second).C:
					log.Println("Ticking", k)
					k++
					if k == 2 {
						log.Println("Cancelling ticker", k)
						cf()
					}
				case <-ctx.Done():
					log.Println("Completed inner go routine #2 from simpleHandler()", ctx.Err())
					return
				}
			}
		})

		<-newCtx2.Done()
		log.Println("Completed long-running go routine from simpleHandler()", ctx.Err())
	})
}

func main() {
	http.HandleFunc("/hello", simpleHandler)
	svr := &http.Server{
		Addr:    ":8081",
		Handler: http.HandlerFunc(simpleHandler),
	}

	graceful.Go(context.Background(),
		func(ctx context.Context) {
			log.Println("Starting server")
			err := svr.ListenAndServe()
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				panic(err)
			}
		})
	graceful.Go(context.Background(),
		func(ctx context.Context) {
			<-ctx.Done()
			log.Println("Shutting down server")
			if err := svr.Shutdown(context.Background()); err != nil {
				log.Printf("Shutdown error: %v", err)
			}
		},
	)

	// Example go routine tracking.
	graceful.Go(context.Background(),
		func(ctx context.Context) {
			deadline := time.NewTimer(time.Second * 5)
			<-deadline.C
			log.Println("Completed running go routine from main().", ctx.Err())
		})

	// Wait for signals from os.
	graceful.Wait()
}
