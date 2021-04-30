package orchestrator

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var sc *sync.Cond
var handlingCount int32

func init() {
	m := &sync.Mutex{}
	sc = sync.NewCond(m)
	atomic.StoreInt32(&handlingCount, 0)

}
func Throttling(requestID string) {
	sc.L.Lock()
	fmt.Printf("c %v \n", atomic.LoadInt32(&handlingCount))
	// 符合就等
	for atomic.LoadInt32(&handlingCount) >= 3 {
		fmt.Println("req" + requestID + "is waiting")
		sc.Wait()
	}
	atomic.AddInt32(&handlingCount, 1)
	sc.L.Unlock()
	// 換我了
	time.Sleep(time.Second)
	atomic.AddInt32(&handlingCount, -1)
}
