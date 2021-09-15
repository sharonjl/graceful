package graceful

import (
	"context"
	"os"
	"os/signal"
)

type Manager struct {
	c        Uint64
	graceCtx context.Context
	cancel   context.CancelFunc
}

// New creates a manager with the given context.
func New() *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	return &Manager{graceCtx: ctx, cancel: cancel}
}

// Go wraps and executes the given function w with a background context.
func (m *Manager) Go(w func(ctx context.Context)) {
	m.GoCtx(context.Background(), w)
}

// GoCtx wraps and executes the given function w with the context given
// by the caller.
func (m *Manager) GoCtx(ctx context.Context, w func(ctx context.Context)) {
	m.c.Inc()
	go func() {
		ctx, cancel := context.WithCancel(ctx)
		defer func() {
			cancel()
			m.c.Dec()
		}()
		go func() {
			select {
			case <-m.graceCtx.Done():
				cancel()
			case <-ctx.Done():
				// do nothing
			}
		}()
		w(ctx)
	}()
}

var sig = make(chan os.Signal, 1)

// Wait listens for the notification signals from the os. When a signal
// is received the main context's cancel is called forcing a cancellation
// of all contexts.
//
// Waits on os.Interrupt and os.Kill when sigs argument is omitted.
func (m *Manager) Wait(sigs ...os.Signal) {
	if len(sigs) == 0 {
		sigs = append(sigs, os.Interrupt, os.Kill)
	}
	signal.Notify(sig, sigs...)

	<-sig
	m.cancel()
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

// Go wraps and executes the given function w with a background context.
func Go(w func(ctx context.Context)) {
	mgr.Go(w)
}

// GoCtx wraps and executes the given function w with the context given
// by the caller.
func GoCtx(ctx context.Context, w func(ctx context.Context)) {
	mgr.GoCtx(ctx, w)
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
