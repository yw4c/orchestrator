package orchestrator

import (
	"github.com/magiconair/properties/assert"
	"orchestrator/pkg/pkgerror"
	"strconv"
	"sync"
	"testing"
	"time"
)

type mockAsyncDTO struct {
	*AsyncFlowContext
	foo string
}

func TestReqWait_Wait(t *testing.T) {

	var wg sync.WaitGroup
	var concurrency = 100
	var timeout = 10*time.Second
	var finishedCount = 0
	wg.Add(concurrency)

	for i:= 1;i<=concurrency ;i++ {

		go func(i int) {
			reqId := "req"+strconv.Itoa(i)

			// start waiting
			data, err := Wait(reqId, timeout, func() {
				wg.Done()
			})

			if err != nil {
				t.Log(err.Error())
				t.Fail()
			} else {
				if resp, ok := data.(*mockAsyncDTO); ok {
					assert.Equal(t, resp.foo, "bar")
				} else {
					t.Log("can not parse to mockAsyncDTO")
					t.Fail()
				}
			}
			finishedCount++

		}(i)
	}
	wg.Wait()
	t.Log("registered")

	for finishedCount < concurrency {
		t.Log(finishedCount)
		for i:=1; i<=concurrency; i++ {
			reqId := "req"+strconv.Itoa(i)

			// mock Msg
			dao := &mockAsyncDTO{
				AsyncFlowContext: &AsyncFlowContext{
					RollbackTopic: "rollback",
					RequestID:     reqId,
					CurrentIndex:  0,
					Topics:        []Topic{"foo"},
				},
				foo: "bar",
			}
			TaskFinished("req"+strconv.Itoa(i), dao)
		}
		time.Sleep(time.Second)
	}

}

func TestReqWait_Wait_Timeout(t *testing.T) {

	var wg sync.WaitGroup
	var errCount int
	var timeout = 1*time.Second
	var mockTimeout = 2*time.Second
	wg.Add(10)


	for i:= 1;i<=10 ;i++ {

		go func(i int) {
			reqId := "req"+strconv.Itoa(i)
			//  start waiting
			_, err := Wait(reqId, timeout, func() {
				wg.Done()
			})
			if err != nil {
				grpcErr := pkgerror.SetGRPCErrorResp("req"+strconv.Itoa(i), err)
				t.Log(grpcErr.Error())
				errCount++
			}

		}(i)
	}
	wg.Wait()

	for i:=1; i<=5; i++ {
		err := TaskFinished("req"+strconv.Itoa(i), &mockAsyncDTO{})
		if err != nil {
			t.Log(err)
			t.Fail()
		}
	}
	// mock timeout
	time.Sleep(mockTimeout)
	for i:=6; i<=10; i++ {
		//finishedRequestsWatcher.Broadcast()
		TaskFinished("req"+strconv.Itoa(i), &mockAsyncDTO{})
	}


	t.Log("timeout request count: ", errCount)
	assert.Equal(t, errCount, 5)
}