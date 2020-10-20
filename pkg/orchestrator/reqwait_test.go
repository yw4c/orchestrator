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
	var concurrency = 10
	var timeout = 1*time.Second
	wg.Add(concurrency)

	for i:= 1;i<=concurrency ;i++ {

		go func(i int) {
			reqId := "req"+strconv.Itoa(i)

			//  start waiting
			data, err := Wait(reqId, timeout)

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

			wg.Done()
		}(i)
	}

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

	wg.Wait()
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
			_, err := Wait(reqId, timeout)
			if err != nil {
				grpcErr := pkgerror.SetGRPCErrorResp("req"+strconv.Itoa(i), err)
				t.Log(grpcErr.Error())
				errCount++
			}
			wg.Done()
		}(i)
	}

	for i:=1; i<=5; i++ {
		TaskFinished("req"+strconv.Itoa(i), &mockAsyncDTO{})
	}
	// mock timeout
	time.Sleep(mockTimeout)
	for i:=6; i<=10; i++ {
		TaskFinished("req"+strconv.Itoa(i), &mockAsyncDTO{})
	}

	wg.Wait()
	t.Log("timeout request count: ", errCount)
	assert.Equal(t, errCount, 5)
}