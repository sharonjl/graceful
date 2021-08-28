package graceful

import (
	"context"
	"os"
	"os/signal"
)

type Manager struct {
	c   atomicInt64
	ctx context.Context
	cf  context.CancelFunc
}

// New creates a manager with the given context.
func New() *Manager {
	ctx, cf := context.WithCancel(context.Background())
	return &Manager{
		ctx: ctx,
		cf:  cf,
	}
}

// Go wraps and executes the given function w with a cancellable context.
func (m *Manager) Go(ctx context.Context, w func(ctx context.Context)) {
	m.c.inc()
	go func() {
		wCtx, cancel := context.WithCancel(ctx)
		defer func() {
			cancel()
			m.c.dec()
		}()

		go func() {
			select {
			case <-m.ctx.Done():
				cancel()
			}
		}()
		w(wCtx)
	}()
}

var (
	defaultWaitSignals = []os.Signal{os.Interrupt, os.Kill}
	sig                = make(chan os.Signal, 1)
)

// Wait listens for the notification signals from the os, when a signal
// is received the manager's context is cancelled.
//
// When the sigs argument is omitted, we wait on the signals defined
// in defaultWaitSignals.
func (m *Manager) Wait(sigs ...os.Signal) {
	if len(sigs) == 0 {
		sigs = append(sigs, defaultWaitSignals...)
	}
	signal.Notify(sig, sigs...)

	<-sig
	m.cf()
	for !m.c.isZero() {
		// wait for routines to exit
	}
}

// Count returns the number of routines Manager is tracking right now.
func (m *Manager) Count() int64 {
	return m.c.get()
}

// mgr Package level instance of Manager.
var mgr = New()

// Go wraps and executes the given function w with a cancellable context.
// When the manager's context is cancelled, the cancel function for w is
// called.
func Go(ctx context.Context, w func(ctx context.Context)) {
	mgr.Go(ctx, w)
}

// Wait listens for the provided notification signals from the os.
// When the signal is received the manager's context is cancelled.
//
// When the sigs argument is omitted, we wait on the signals defined
// in defaultWaitSignals.
func Wait(sigs ...os.Signal) {
	mgr.Wait(sigs...)
}

// Count returns the number of routines graceful is tracking right now.
func Count() int64 {
	return mgr.Count()
}
