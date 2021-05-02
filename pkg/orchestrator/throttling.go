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

// pendingRequestMap : [string(request-id)]pendingRequest
var pendingRequestMap = &sync.Map{}

type pendingRequest struct {
	respChan chan interface{}
}

type TaskDone func()

func init() {
	go func() {
		for atomic.LoadInt32(&handlingCount) >= int32(config.GetConfigInstance().Throttling.Concurrency) {
			pendingRequestMap.LoadAndDelete()
		}
	}()
}

func Throttling(requestID string, timeout time.Duration) (err error, callback TaskDone) {

	pending := &pendingRequest{
		respChan: make(chan interface{}),
	}
	atomic.AddInt32(&handlingCount, 1)

	// waiting if present concurrency count is over than passed concurrency count
	log.Debug().Str("request ID", requestID).Msg("Waiting Task ")
	pendingRequestMap.Store(requestID, pending)
	<-pending.respChan

	return nil, func() {
		atomic.AddInt32(&handlingCount, -1)

		log.Debug().Str("request ID", requestID).Msg("Done Task ")
	}
}

// func TaskDone(requestID string, resp interface{}) error {

// 	if pending, ok := pendingRequestMap.Load(requestID); ok {
// 		if d, ok := pending.(*pendingRequest); ok {
// 			d.respChan <- resp
// 		} else {
// 			return eris.Wrap(pkgerror.ErrInternalError, "convert failed ")
// 		}
// 	} else {
// 		return eris.Wrap(pkgerror.ErrInternalError, "pending request has not been registered")
// 	}
// 	return nil

// }
