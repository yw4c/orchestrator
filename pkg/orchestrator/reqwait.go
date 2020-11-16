package orchestrator

import (
	"fmt"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog/log"
	"orchestrator/pkg/pkgerror"
	"sync"
	"time"
)

var finishedRequests = &sync.Map{}
var finishedRequestsWatcher = sync.NewCond(new(sync.Mutex))

func Wait(requestID string, timeout time.Duration) (dto IAsyncFlowContext, err error) {

	timestamp := time.Now()
	var data interface{}
	isReqFinished := func()bool {
		var ok bool
		data,ok = finishedRequests.Load(requestID)
		return ok
	}


	// if still pending
	for !isReqFinished() {
		fmt.Println(time.Now().Unix(), timestamp.Add(timeout).Unix())

		// cancel while timeout
		if time.Now().After(timestamp.Add(timeout)) {
			err = eris.Wrap(pkgerror.ErrTimeout, "async flow handles it too long, let us cancel it")
			return
		}

		// keep waiting
		log.Debug().Str("request ID", requestID).Msg("Waiting Task ")

		// sync.Map doesn't need to be locked
		finishedRequestsWatcher.L.Lock()
		finishedRequestsWatcher.Wait()
		finishedRequestsWatcher.L.Unlock()
	}

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

	finishedRequests.Delete(requestID)
	return


}

func TaskFinished(requestID string, resp interface{}) {
	finishedRequests.Store(requestID, resp)
	finishedRequestsWatcher.Broadcast()
}