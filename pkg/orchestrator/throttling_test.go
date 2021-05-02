package orchestrator

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestReqWait_Throttling(t *testing.T) {

	var wg sync.WaitGroup
	var concurrency = 10
	var timeout = 10 * time.Second
	var finishedCount = 0
	wg.Add(concurrency)

	for i := 1; i <= concurrency; i++ {

		go func(i int) {
			reqId := "req" + strconv.Itoa(i)

			// start waiting
			_, taskDone := Throttling(reqId, timeout)
			// do stuff
			taskDone()
			t.Log(reqId + "done")
			finishedCount++

		}(i)
	}
	wg.Wait()

}

// func TestReqWait_Throttling_Timeout(t *testing.T) {

// 	var wg sync.WaitGroup
// 	var errCount int
// 	var timeout = 1 * time.Second
// 	var mockTimeout = 2 * time.Second
// 	wg.Add(10)

// 	for i := 1; i <= 10; i++ {

// 		go func(i int) {
// 			reqId := "req" + strconv.Itoa(i)
// 			//  start waiting
// 			_, err := Wait(reqId, timeout, func() {
// 				wg.Done()
// 			})
// 			if err != nil {
// 				grpcErr := pkgerror.SetGRPCErrorResp("req"+strconv.Itoa(i), err)
// 				t.Log(grpcErr.Error())
// 				errCount++
// 			}

// 		}(i)
// 	}
// 	wg.Wait()

// 	for i := 1; i <= 5; i++ {
// 		err := TaskDone("req"+strconv.Itoa(i), &mockAsyncDTO{})
// 		if err != nil {
// 			t.Log(err)
// 			t.Fail()
// 		}
// 	}
// 	// mock timeout
// 	time.Sleep(mockTimeout)
// 	for i := 6; i <= 10; i++ {
// 		//finishedRequestsWatcher.Broadcast()
// 		TaskDone("req"+strconv.Itoa(i), &mockAsyncDTO{})
// 	}

// 	t.Log("timeout request count: ", errCount)
// 	assert.Equal(t, errCount, 5)
// }
