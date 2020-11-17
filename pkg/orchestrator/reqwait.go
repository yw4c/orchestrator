package orchestrator

import (
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog/log"
	"orchestrator/pkg/pkgerror"
	"sync"
	"time"
)

// pendingRequestMap : [string(request-id)]pendingRequest
var pendingRequestMap = &sync.Map{}

type pendingRequest struct {
	//requestID string
	//context context.Context
	// respChan: IAsyncFlowContext or error
	respChan chan interface{}
}

func Wait(requestID string, timeout time.Duration, startFunc func()) (dto IAsyncFlowContext, err error) {

	pending := &pendingRequest{
		respChan: make(chan interface{}),
	}

	pendingRequestMap.Store(requestID, pending)
	startFunc()


	// keep waiting
	log.Debug().Str("request ID", requestID).Msg("Waiting Task ")

	select {
		case <- time.Tick(timeout):
			err = eris.Wrap(pkgerror.ErrTimeout, "async flow handles it too long, let us cancel it")
			pendingRequestMap.Delete(requestID)
			return
		case data := <- pending.respChan:
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

func TaskFinished(requestID string, resp interface{}) error {

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