package graceful

import "sync/atomic"

type Uint64 struct {
	v uint64
}

// NewUint64 creates a new Uint64.
func NewUint64(val uint64) *Uint64 {
	return &Uint64{v: val}
}

// Load atomically loads the wrapped value.
func (i *Uint64) Load() uint64 {
	return atomic.LoadUint64(&i.v)
}

// Add atomically adds to the wrapped uint64 and returns the new value.
func (i *Uint64) Add(delta uint64) uint64 {
	return atomic.AddUint64(&i.v, delta)
}

// Sub atomically subtracts from the wrapped uint64 and returns the new value.
func (i *Uint64) Sub(delta uint64) uint64 {
	return atomic.AddUint64(&i.v, ^(delta - 1))
}

// Inc atomically increments the wrapped uint64 and returns the new value.
func (i *Uint64) Inc() uint64 {
	return i.Add(1)
}

// Dec atomically decrements the wrapped uint64 and returns the new value.
func (i *Uint64) Dec() uint64 {
	return i.Sub(1)
}

// CAS is an atomic compare-and-swap.
func (i *Uint64) CAS(old, new uint64) (swapped bool) {
	return atomic.CompareAndSwapUint64(&i.v, old, new)
}
