package graceful

type RoutineRunner struct {
	c atomicInt64
	s func(f func())
}

func (g *RoutineRunner) Terminate() {
	for !g.c.isZero() {
		// wait for routines to exit
	}
}

func (g *RoutineRunner) Run(f func()) {
	g.c.inc()
	go func() {
		defer g.c.dec()

		if g.s == nil {
			f()
		} else {
			g.s(f)
		}
	}()
}

// GoRoutineTerminator returns the underlying DefaultRunner for go routines
// wrapped in a Terminator for graceful termination.
func GoRoutineTerminator() Terminator {
	return &DefaultRunner
}
