package node

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"orchestrator/pkg/orchestrator"
)

func CancelBooking() orchestrator.RollbackNode {
	return func(topic orchestrator.Topic, data []byte) {

		msg := &orchestrator.RollbackMsg{}
		if err := json.Unmarshal(data, msg); err != nil {
			log.Panic().Str("topic", string(topic)).Str("data", string(data)).Msg("Failed to unmarshal")
			return
		}

		log.Info().Msg("Assume Invoking PaymentService.CancelOrder Success !")
		log.Info().Msg("Assume Invoking PaymentService.CancelPayment Success !")

		log.Info().
			Str("topic", string(topic)).
			Str("Message", string(data)).
			Str("requestID", msg.RequestID).
			Msg("Cancelling Booking, Rollback success !")

	}
}
