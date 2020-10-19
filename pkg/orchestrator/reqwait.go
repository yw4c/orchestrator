package orchestrator

import (
	"context"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

var FinishedRequests = sync.Map{}

func Wait(requestID string, timeout time.Duration, ctx context.Context) {
	log.Debug().Str("request ID", requestID).Msg("Waiting Task Finish")


	for {
		//fmt.Println(FinishedRequests)
		if _, ok := FinishedRequests.Load(requestID) ; ok {
			log.Debug().
				Str("request ID", requestID).
				Msg("Finished Task")
			FinishedRequests.Delete(requestID)
			return
		}
	}

}

func TaskFinished(requestID string) {
	FinishedRequests.Store(requestID, 1)
}