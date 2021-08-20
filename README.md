# graceful
graceful is a resource termination library to smoothly clean up resources on term signals.

# example
```go
package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sharonjl/graceful"
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
	
	graceful.Go(func() {
		for range time.NewTimer(time.Second * 5).C {
			// do nothing
		}
	})
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
```