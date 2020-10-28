package orchestrator

import (
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog/log"
	"orchestrator/pkg/pkgerror"
	"sync"
	"time"
)

var FinishedRequests = &sync.Map{}

func Wait(requestID string, timeout time.Duration) (dto IAsyncFlowContext, err error) {
	log.Debug().Str("request ID", requestID).Msg("Waiting Task ")

	timestamp := time.Now()

	for {
		// if waited requestID got into the map
		if data, ok := FinishedRequests.Load(requestID) ; ok {

			// convert to IAsyncFlowContext
			if d, ok := data.(IAsyncFlowContext); ok {
				dto = d
			} else if e, ok := data.(error); ok { // convert to error
				err = e
			} else {
				err = eris.Wrap(pkgerror.ErrInternalError, "convert to IAsyncFlowContext failed ")
			}

			log.Debug().
				Str("request ID", requestID).
				Msg("Finished Task")

			FinishedRequests.Delete(requestID)
			return
		}

		// cancel while timeout
		if time.Now().After(timestamp.Add(timeout)) {
			err = eris.Wrap(pkgerror.ErrTimeout, "async flow handles it too long, let us cancel it")
			return
		}

		time.Sleep(100*time.Millisecond)
	}

}

func TaskFinished(requestID string, resp interface{}) {
	FinishedRequests.Store(requestID, resp)
}