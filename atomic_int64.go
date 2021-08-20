package graceful

import "sync/atomic"

type atomicInt64 int64

func (c *atomicInt64) inc() int64 {
	return atomic.AddInt64((*int64)(c), 1)
}

func (c *atomicInt64) dec() int64 {
	return atomic.AddInt64((*int64)(c), -1)
}

func (c *atomicInt64) isZero() bool {
	return atomic.CompareAndSwapInt64((*int64)(c), 0, 0)
}
