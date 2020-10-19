package orchestrator

import (
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog/log"
	"orchestrator/pkg/pkgerror"
	"sync"
	"time"
)

var FinishedRequests = &sync.Map{}

func Wait(requestID string, timeout time.Duration) error {
	log.Debug().Str("request ID", requestID).Msg("Waiting Task Finish")

	timestamp := time.Now()

	for {
		// return if we would got the request finished
		if _, ok := FinishedRequests.Load(requestID) ; ok {
			log.Debug().
				Str("request ID", requestID).
				Msg("Finished Task")
			FinishedRequests.Delete(requestID)
			return nil
		}

		// cancel while timeout
		if time.Now().After(timestamp.Add(timeout)) {
			return eris.Wrap(pkgerror.ErrTimeout, "async flow handles it too long, let us cancel it")
		}
	}

}

func TaskFinished(requestID string) {
	FinishedRequests.Store(requestID, 1)
}