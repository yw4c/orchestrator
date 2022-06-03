package node

import (
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog/log"
	"orchestrator/pb"
	"orchestrator/pkg/orchestrator"
	"orchestrator/pkg/pkgerror"
)

func CreateOrderSyncNode() orchestrator.SyncNode {
	return func(requestID string, ctx *orchestrator.Context) error {

		// load context
		request := &pb.BookingRequest{}
		ctx.LoadPbReq(request)
		resp := &pb.BookingSyncResponse{}
		ctx.LoadPbResp(resp)

		log.Info().Msg("Assume Invoking Order Service Success !")

		// Set response in context
		resp.OrderID = 1
		ctx.SetPbResp(resp)

		return nil
	}

}

func CreatePaymentSyncNode() orchestrator.SyncNode {
	return func(requestID string, ctx *orchestrator.Context) error {

		// load request payload
		request := &pb.BookingRequest{}
		ctx.LoadPbReq(request)
		resp := &pb.BookingSyncResponse{}
		ctx.LoadPbResp(resp)

		// Fault Injection for observing rollback
		if request.FaultInject {
			log.Error().Msg("Assume Invoking Payment Service Fail !")
			return eris.Wrap(pkgerror.ErrInvalidInput, "Fault Injection")
		}

		log.Info().Msg("Assume Invoking Payment Service Success !")
		// Set response in context
		resp.PaymentID = 2
		ctx.SetPbResp(resp)
		return nil
	}
}
