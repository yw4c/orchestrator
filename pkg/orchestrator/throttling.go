package orchestrator

import (
	"orchestrator/config"
	"orchestrator/pkg/pkgerror"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rotisserie/eris"
	"github.com/rs/zerolog/log"
)

// count of handling request
var handlingCount int32

// pendingRequestMap : [string(request-id)]pendingRequest
var pendingRequestMap = &sync.Map{}

type pendingRequest struct {
	respChan chan interface{}
}

func Throttling(requestID string, timeout time.Duration) (err error) {

	pending := &pendingRequest{
		respChan: make(chan interface{}),
	}
	pendingRequestMap.Store(requestID, pending)

	// keep waiting
	log.Debug().Str("request ID", requestID).Msg("Waiting Task ")

	for atomic.LoadInt32(&handlingCount) >= int32(config.GetConfigInstance().Throttling.Concurrency) {
		select {
		case <-time.Tick(timeout):
			err = eris.Wrap(pkgerror.ErrTimeout, "async flow handles it too long, let us cancel it")
			pendingRequestMap.Delete(requestID)
			return
		case data := <-pending.respChan:
			// finished waiting
			log.Debug().Str("request ID", requestID).Msg("Finished Task")

			// convert to IAsyncFlowContext
			if d, ok := data.(IAsyncFlowContext); ok {
				dto = d
			} else if e, ok := data.(error); ok { // convert to error
				err = e
			} else {
				err = eris.Wrap(pkgerror.ErrInternalError, "convert to IAsyncFlowContext failed ")
			}
			pendingRequestMap.Delete(requestID)
			return
		}
	}

}

func TaskDone(requestID string, resp interface{}) error {

	if pending, ok := pendingRequestMap.Load(requestID); ok {
		if d, ok := pending.(*pendingRequest); ok {
			d.respChan <- resp
		} else {
			return eris.Wrap(pkgerror.ErrInternalError, "convert failed ")
		}
	} else {
		return eris.Wrap(pkgerror.ErrInternalError, "pending request has not been registered")
	}
	return nil

}
