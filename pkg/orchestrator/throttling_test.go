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
			fmt.Printf("req %v is done \n", i)
		}(i)
	}

	// access 2
	time.Sleep(2 * time.Second)
	fmt.Println(atomic.LoadInt32(&handlingCount))
	sc.Broadcast()

	// access 3
	time.Sleep(2 * time.Second)
	fmt.Println(atomic.LoadInt32(&handlingCount))
	sc.Broadcast()

	select {}
}
