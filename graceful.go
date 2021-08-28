package graceful

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// DefaultWaitSignals a list of os.Signals to be notified on for the Wait operation.
var DefaultWaitSignals = []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL}

var sig = make(chan os.Signal, 1)

type Manager struct {
	c        atomicInt64
	t        atomicInt64
	lifeline chan struct{}
	mu       sync.Mutex
	m        map[int64]context.CancelFunc
}

// New creates a manager with the given context.
func New() *Manager {
	return &Manager{
		lifeline: make(chan struct{}),
		m:        make(map[int64]context.CancelFunc),
		mu:       sync.Mutex{},
	}
}

// Go wraps and executes the given function w with a lifeline context.
func (m *Manager) Go(ctx context.Context, w func(ctx context.Context)) {
	m.c.inc()
	n := m.t.inc()
	go func() {
		defer func() {
			if cf, ok := m.m[n]; ok && cf != nil {
				cf()
			}
			m.mu.Lock()
			delete(m.m, n)
			m.mu.Unlock()
			m.c.dec()
		}()
		cancelCtx, cancelFunc := context.WithCancel(ctx)

		m.mu.Lock()
		m.m[n] = cancelFunc
		m.mu.Unlock()
		w(cancelCtx)
	}()
}

// Wait listens for the provided notification signals from the os.
// When the signal is received the manager's lifeline channel is closed.
//
// When the sigs argument is omitted, we wait on the signals defined
// in DefaultWaitSignals.
func (m *Manager) Wait(sigs ...os.Signal) {
	if len(sigs) == 0 {
		sigs = append(sigs, DefaultWaitSignals...)
	}
	signal.Notify(sig, sigs...)

	<-sig
	m.mu.Lock()
	for k, cf := range m.m {
		cf()
		m.m[k] = nil
	}
	m.mu.Unlock()
	for !m.c.isZero() {
		// wait for routines to exit
	}
}

// Count returns the number of routines Manager is tracking right now.
func (m *Manager) Count() int64 {
	return m.c.get()
}

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
// in DefaultWaitSignals.
func Wait(sigs ...os.Signal) {
	mgr.Wait(sigs...)
}

// Count returns the number of routines graceful is tracking right now.
func Count() int64 {
	return mgr.Count()
}
