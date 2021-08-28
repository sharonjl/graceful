package graceful

import "sync/atomic"

type atomicInt64 int64

func (i *atomicInt64) get() int64 {
	return atomic.LoadInt64((*int64)(i))
}

func (i *atomicInt64) inc() int64 {
	return atomic.AddInt64((*int64)(i), 1)
}

func (i *atomicInt64) dec() int64 {
	return atomic.AddInt64((*int64)(i), -1)
}

func (i *atomicInt64) isZero() bool {
	return atomic.CompareAndSwapInt64((*int64)(i), 0, 0)
}
