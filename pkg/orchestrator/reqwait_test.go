package orchestrator

import (
	"fmt"
	"orchestrator/pkg/pkgerror"
	"strconv"
	"testing"
	"time"
)

func TestReqWait_Wait(t *testing.T) {

	for i:= 1;i<10 ;i++ {
		go func(i int) {
			err := Wait("req"+strconv.Itoa(i), 5*time.Second)
			grpcErr := pkgerror.SetGRPCErrorResp("req"+strconv.Itoa(i), err)
			fmt.Println(grpcErr)
		}(i)
	}
	time.Sleep(10*time.Second)

	for i:=1; i<12; i++ {
		TaskFinished("req"+strconv.Itoa(i))
		time.Sleep(time.Second)
	}

	t.Log("finish")

	select {

	}
}
