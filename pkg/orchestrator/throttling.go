package orchestrator

import (
	"orchestrator/config"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog/log"
)

// count of handling request
var handlingCount int32

var pendingRequestChan chan *pendingRequest
var throttlingMu sync.Mutex

type pendingRequest struct {
	requestID string
	respChan  chan interface{}
}

// callback after task done
type TaskDone func()

func init() {
	pendingRequestChan = make(chan *pendingRequest)

	// todo singlton
	go func() {
		for {
			select {
			// hey buddy, get off from queue!
			case pending := <-pendingRequestChan:
				log.Debug().Str("request ID", pending.requestID).Msg("Received ")

				// go to wait your self
				go func(*pendingRequest) {

					for {
						throttlingMu.Lock()
						if atomic.LoadInt32(&handlingCount) < int32(config.GetConfigInstance().Throttling.Concurrency) {
							log.Debug().
								Str("id", pending.requestID).
								Int32("handling", atomic.LoadInt32(&handlingCount)).
								Msg("Check count")
							atomic.AddInt32(&handlingCount, 1)
						} else {
							throttlingMu.Unlock()
							continue
						}
						throttlingMu.Unlock()
						pending.respChan <- 1
						return
					}

				}(pending)

			}
		}
	}()
}

func Throttling(requestID string, timeout time.Duration) (err error, callback TaskDone) {

	pending := &pendingRequest{
		requestID: requestID,
		respChan:  make(chan interface{}),
	}

	// send msg to queue that i'm waitting !
	pendingRequestChan <- pending

	// wait
	<-pending.respChan

	return nil, func() {
		atomic.AddInt32(&handlingCount, -1)
		log.Debug().Str("request ID", requestID).Int32("handling", atomic.LoadInt32(&handlingCount)).Msg("Done Task ")
	}
}
