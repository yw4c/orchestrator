package orchestrator

import (
	"sync"
	"sync/atomic"
)

var sc *sync.Cond

// count of handling request
var handlingCount int32

// if there are waiting request, keep broadcasting
var waitingChan chan int

func init() {
	m := &sync.Mutex{}
	sc = sync.NewCond(m)
	atomic.StoreInt32(&handlingCount, 0)
	waitingChan = make(chan int)

	go func() {
		for range waitingChan {
			sc.Broadcast()
		}
	}()
}

func Throttling(requestID string) {
	sc.L.Lock()
	// 處理中的請求達上限就等
	for atomic.LoadInt32(&handlingCount) >= 3 {
		waitingChan <- 1
		sc.Wait()
	}
	atomic.AddInt32(&handlingCount, 1)
	sc.L.Unlock()
}
