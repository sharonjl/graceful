package graceful

import (
	"context"
	"errors"
	"sync"
)

type lifelineCtx struct {
	context.Context

	// lifeline is a control channel which when closed
	// will close/cancel all child contexts.
	lifeline chan struct{}

	mu   sync.Mutex    // protects the following fields
	done chan struct{} // context's done channel
	err  error         // context's err channel
}

func WithLifeline(parent context.Context, lifeline chan struct{}) context.Context {
	return &lifelineCtx{Context: parent, lifeline: lifeline, mu: sync.Mutex{}}
}

var ErrLifelineTerminated = errors.New("context lifeline terminated")

func (c *lifelineCtx) Done() <-chan struct{} {
	c.mu.Lock()
	if c.done == nil {
		c.done = make(chan struct{})
		go func() {
			// Close this context's done channel when either the lifeline
			// or the parent context's done is closed, this allows cancel
			// to propagate through the context tree.
			select {
			case <-c.Context.Done():
				c.mu.Lock()
				c.err = c.Context.Err()
				c.mu.Unlock()
			case <-c.lifeline:
				c.mu.Lock()
				c.err = ErrLifelineTerminated
				c.mu.Unlock()
			}
			close(c.done)
		}()
	}
	d := c.done
	c.mu.Unlock()
	return d
}

func (c *lifelineCtx) Err() error {
	c.mu.Lock()
	err := c.err
	c.mu.Unlock()
	return err
}

func (c *lifelineCtx) String() string {
	return "graceful.WithLifeline"
}
