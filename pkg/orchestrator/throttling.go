package orchestrator

import (
	"fmt"
	"orchestrator/config"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog/log"
)

var sc *sync.Cond

// count of handling request
var handlingCount int32

// if there are waiting request, keep broadcasting
var waitingChan chan int

func init() {

	if !config.GetConfigInstance().Throttling.Enable {
		return
	}
	m := &sync.Mutex{}
	sc = sync.NewCond(m)
	atomic.StoreInt32(&handlingCount, 0)
	waitingChan = make(chan int)

	go func() {
		for {
			<-waitingChan
			sc.Broadcast()
			fmt.Println("Broadcast")
			time.Sleep(time.Second)
		}
	}()
	log.Info().Msg("throttling registered!")
}

func throttling(id string) {
	if !config.GetConfigInstance().Throttling.Enable {
		return
	}
	sc.L.Lock()
	// 處理中的請求達上限就等
	for atomic.LoadInt32(&handlingCount) >= int32(config.GetConfigInstance().Throttling.Concurrency) {
		fmt.Println(id + "waiting")
		waitingChan <- 1
		fmt.Println(id + "waiting2")
		time.Sleep(time.Second)
		sc.Wait()
	}
	atomic.AddInt32(&handlingCount, 1)
	sc.L.Unlock()
}
