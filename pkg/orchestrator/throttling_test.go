package orchestrator

import (
	"fmt"
	"strconv"
	"sync/atomic"
	"testing"
	"time"
)

func TestThrottling(t *testing.T) {

	for i := 0; i < 10; i++ {
		go func(i int) {
			Throttling(strconv.Itoa(i))
			// do your task
			time.Sleep(2 * time.Second)
			fmt.Printf("req %v is done \n", i)
			defer atomic.AddInt32(&handlingCount, -1)
		}(i)
	}
	select {}
}
