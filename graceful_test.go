package graceful

import (
	"testing"
	"time"
)

func Test_GoRoutines(t *testing.T) {
	Go(func() {
		t.Log("routine #1: waiting 2s")
		<-time.After(time.Second * 2)
		t.Log("routine #1: done")
		Go(func() {
			t.Log("routine #2: waiting 1s")
			<-time.After(time.Second * 1)
			t.Log("routine #2: done")
			Go(func() {
				t.Log("routine #3: waiting 2s")
				<-time.After(time.Second * 2)
				t.Log("routine #3: resuming, will wait again for 3s")
				<-time.After(time.Second * 3)
				t.Log("routine #3: done")
			})
		})
		Go(func() {
			t.Log("routine #4: waiting 1s")
			<-time.After(time.Second * 1)
			t.Log("routine #4: done")
		})
	})

	Go(func() {
		t.Log("routine #5: waiting 2s")
		<-time.After(time.Second * 2)
		t.Log("routine #5: done")
	})

	Go(func() {
		t.Log("routine #6: waiting 3s")
		<-time.After(time.Second * 3)
		t.Log("routine #6: done")
	})

	close(sig)
	Wait()
}
