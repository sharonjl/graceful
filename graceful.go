package graceful

import (
	"context"
	"os"
	"os/signal"
)

type Manager struct {
	c        Uint64
	t        Uint64
	lifeline chan struct{}
}

// New creates a manager with the given context.
func New() *Manager {
	return &Manager{lifeline: make(chan struct{})}
}

// Go wraps and executes the given function w with a cancellable context.
func (m *Manager) Go(ctx context.Context, w func(ctx context.Context)) {
	m.c.Inc()
	go func() {
		defer m.c.Dec()
		w(WithLifeline(ctx, m.lifeline))
	}()
}

var sig = make(chan os.Signal, 1)

// Wait listens for the notification signals from the os. When a signal
// is received context.CancelFunc is called for contexts (go routines)
// being tracked.
//
// Waits on os.Interrupt and os.Kill when sigs argument is omitted.
func (m *Manager) Wait(sigs ...os.Signal) {
	if len(sigs) == 0 {
		sigs = append(sigs, os.Interrupt, os.Kill)
	}
	signal.Notify(sig, sigs...)

	<-sig
	close(m.lifeline)
	for !m.c.CAS(0, 0) {
		// wait for routines to exit
	}
}

// Count returns the number of routines Manager is tracking right now.
func (m *Manager) Count() uint64 {
	return m.c.Load()
}

// mgr Package level instance of Manager.
var mgr = New()

// Go wraps and executes the given function w with a cancellable context.
// When the manager's context is cancelled, the cancel function for w is
// called.
func Go(ctx context.Context, w func(ctx context.Context)) {
	mgr.Go(ctx, w)
}

// Wait listens for the notification signals from the os. When a signal
// is received context.CancelFunc is called for contexts (go routines)
// being tracked.
//
// Waits on os.Interrupt and os.Kill when sigs argument is omitted.
func Wait(sigs ...os.Signal) {
	mgr.Wait(sigs...)
}

// Count returns the number of routines graceful is tracking right now.
func Count() uint64 {
	return mgr.Count()
}
