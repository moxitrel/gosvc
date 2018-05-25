package svc

import (
	"testing"
	"time"
)

func TestThunk(t *testing.T) {
	o := NewThunk()
	defer o.Stop()

	o.Do(nil)
	o.Do(func(){
		t1 := time.Now()
		t2 := time.Now()
		t.Logf("%v", t2.Sub(t1))
	})

	time.Sleep(time.Millisecond)
}
