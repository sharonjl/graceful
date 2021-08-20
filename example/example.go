package main

import (
	"context"
	"fmt"
	"net/http"

	"graceful"
)

func simpleHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello world\n")
}

func main() {

	http.HandleFunc("/hello", simpleHandler)
	svr := &http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(simpleHandler),
	}

	graceful.In(httpShutdown(svr), graceful.GoRoutineTerminator())

	go func() {
		err := svr.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()
	graceful.Wait()
}

func httpShutdown(svr *http.Server) graceful.TerminatorFunc {
	return func() {
		if err := svr.Shutdown(context.Background()); err != nil {
			fmt.Println(err)
		}
	}
}
