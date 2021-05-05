package orchestrator

import (
	"orchestrator/config"
	"orchestrator/pkg/pkgerror"
	"sync"
	"sync/atomic"

	"github.com/rotisserie/eris"
	"github.com/rs/zerolog/log"
)

// count of handling request
var handlingCount int32

// queue of pending requests
var pendingRequestChan chan *pendingRequest
var throttlingMu sync.Mutex

type pendingRequest struct {
	requestID string
	passChan  chan interface{}
}

// callback after task done
type TaskDone func()

func init() {

	if !config.GetConfigInstance().Throttling.Enable {
		return
	}

	pendingRequestChan = make(chan *pendingRequest)

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
						// if handling request is fewer than max concurrency of config
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
						pending.passChan <- 1
						return
					}
				}(pending)
			}
		}
	}()
}

func Throttling(requestID string) (err error, callback TaskDone) {
	if !config.GetConfigInstance().Throttling.Enable {
		return eris.Wrap(pkgerror.ErrInternalError, "throttling mode is disable"), nil
	}

	pending := &pendingRequest{
		requestID: requestID,
		passChan:  make(chan interface{}),
	}

	// send msg to queue that i'm waitting !
	pendingRequestChan <- pending

	// wait
	<-pending.passChan

	return nil, func() {
		atomic.AddInt32(&handlingCount, -1)
		log.Debug().Str("request ID", requestID).Int32("handling", atomic.LoadInt32(&handlingCount)).Msg("Done Task ")
	}
}
