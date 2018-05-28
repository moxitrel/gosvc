package svc

import (
	"testing"
	"time"
)

func Test_NewLoopWithNil(t *testing.T) {
	o := NewLoop(nil)
	o.Stop()
}

func Test_Loop(t *testing.T) {
	i := 0
	o := NewLoop(func() {
		i++
	})
	defer o.Stop()

	time.Sleep(time.Millisecond)
	if i == 0 {
		t.Errorf("i == 0, want !0")
	} else {
		t.Logf("i: %v", i)
	}
}

func TestService_Stop(t *testing.T) {
	o := NewLoop(func() {
		time.Sleep(time.Millisecond)
	})
	o.Stop()

	if o.state != STOPPED {
		t.Errorf("o.state != STOPPED, want STOPPED")
	}
}

