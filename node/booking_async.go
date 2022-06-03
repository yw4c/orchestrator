package node

import (
	"encoding/json"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog/log"
	"orchestrator/dto"
	"orchestrator/pkg/orchestrator"
	"orchestrator/pkg/pkgerror"
)

func CreateOrderAsyncNode() orchestrator.AsyncNode {
	return func(topic orchestrator.Topic, data []byte, next orchestrator.Next, rollback orchestrator.Rollback) {

		// Convert MQ Message into struct
		d := &dto.BookingMsgDTO{}

		if err := json.Unmarshal(data, d); err != nil {
			rollback(eris.Wrap(pkgerror.ErrInternalError, "json unmarshal fail"), d)
			return
		}

		// Simulate invoking order service and getting response
		d.OrderID = 11

		next(d)

		log.Info().Msg("CreateOrderAsyncNode finished")
	}
}

func CreatePaymentAsyncNode() orchestrator.AsyncNode {
	return func(topic orchestrator.Topic, data []byte, next orchestrator.Next, rollback orchestrator.Rollback) {

		// Convert Message into struct
		d := &dto.BookingMsgDTO{}

		if err := json.Unmarshal(data, d); err != nil {
			rollback(eris.Wrap(pkgerror.ErrInternalError, "json unmarshal fail"), d)
			return
		}

		// Simulate invoking payment service and getting response
		d.PaymentID = 12

		// Fault Injection for observing rollback
		if d.FaultInject {
			rollback(eris.Wrap(pkgerror.ErrInvalidInput, "this is an mocked invalid error"), d)
			return
		}

		next(d)

		log.Info().Msg("CreatePaymentAsyncNode finished")
	}
}
