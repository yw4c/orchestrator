package orchestrator

import (
	"context"
	"strconv"
	"testing"
	"time"
)

func TestReqWait_Wait(t *testing.T) {

	for i:= 1;i<10 ;i++ {
		go func(i int) {
			Wait("req"+strconv.Itoa(i), 5*time.Second, context.Background())
		}(i)
	}
	time.Sleep(time.Second)

	for i:=1; i<12; i++ {
		TaskFinished("req"+strconv.Itoa(i))
		time.Sleep(time.Second)
	}

	t.Log("finish")

	select {

	}
}
