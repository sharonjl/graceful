package graceful

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Terminator is the interface that wraps the Terminate method.
//
// Terminate instructs the implementing object to shut down,
// close, disposes, or clean up any associated resources.
type Terminator interface {
	Terminate()
}

// TerminatorFunc defines a type on func() as a convenience
type TerminatorFunc func()

// Terminate implements the Terminator interface on TerminatorFunc.
func (f TerminatorFunc) Terminate() {
	f()
}

var (
	tt []Terminator
	mu = sync.Mutex{}
)

// In registers one or more Terminator(s) to execute in order.
func In(f ...Terminator) {
	mu.Lock()
	defer mu.Unlock()

	tt = append(tt, f...)
}

var DefaultRunner = RoutineRunner{}

// Go is used to track go func()
func Go(f func()) {
	DefaultRunner.Run(f)
}

var sig = make(chan os.Signal, 1)

// DefaultWaitSignals a list of os.Signals to be notified on for the Wait operation.
var DefaultWaitSignals = []os.Signal{syscall.SIGINT, syscall.SIGTERM}

// Wait listens for the provided notification signals from the os.
// When the signal is received, all registered terminators are
// executed in sequence.
//
// When the sigs argument is omitted, we wait on the signals defined
// in DefaultWaitSignals. Finally, tracked go routines are terminated
// last if they are not registered for termination using In.
func Wait(sigs ...os.Signal) {
	if len(sigs) == 0 {
		sigs = append(sigs, DefaultWaitSignals...)
	}
	signal.Notify(sig, sigs...)

	<-sig
	didTerminateRoutines := false
	for _, t := range tt {
		if _, ok := t.(*RoutineRunner); ok && !didTerminateRoutines {
			didTerminateRoutines = true
		}
		t.Terminate()
	}
	if !didTerminateRoutines {
		GoRoutineTerminator().Terminate()
	}
}
