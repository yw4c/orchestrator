package orchestrator

import (
	"fmt"
	"strconv"
	"sync/atomic"
	"testing"
	"time"
)

// go test -bench=. -run=^$ throttling_test.go throttling.go
func Test_Throttling(t *testing.T) {
	for i := 0; i < 10; i++ {
		go func(i int) {
			throttling(strconv.Itoa(i))
			// do your task
			time.Sleep(2 * time.Second)
			fmt.Printf("req %v is done \n", i)
			defer atomic.AddInt32(&handlingCount, -1)
		}(i)
	}
	select {}
}
