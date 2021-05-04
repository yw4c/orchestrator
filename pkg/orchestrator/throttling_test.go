package orchestrator

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
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
			// image do stuff
			taskDone()
			wg.Done()
			finishedCount++
		}(i)
	}
	wg.Wait()
	assert.Equal(t, finishedCount, concurrency)

}
