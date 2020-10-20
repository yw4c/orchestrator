package orchestrator

import (
	"fmt"
	"github.com/magiconair/properties/assert"
	"orchestrator/pkg/pkgerror"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestReqWait_Wait(t *testing.T) {

	var wg sync.WaitGroup
	var errCount int
	var timeout = 1*time.Second
	var mockTimeout = 2*time.Second
	wg.Add(10)


	for i:= 1;i<=10 ;i++ {
		go func(i int) {
			//  start waiting
			err := Wait("req"+strconv.Itoa(i), timeout)
			if err != nil {
				grpcErr := pkgerror.SetGRPCErrorResp("req"+strconv.Itoa(i), err)
				fmt.Println(grpcErr)
				errCount++
			}
			wg.Done()
		}(i)
	}

	for i:=1; i<=5; i++ {
		TaskFinished("req"+strconv.Itoa(i))
	}
	// mock timeout
	time.Sleep(mockTimeout)
	for i:=6; i<=10; i++ {
		TaskFinished("req"+strconv.Itoa(i))
	}

	wg.Wait()
	fmt.Println("timeout request count: ", errCount)
	assert.Equal(t, errCount, 5)
}
