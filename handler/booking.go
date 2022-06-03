package handler

import (
	"context"
	"github.com/rotisserie/eris"
	"orchestrator/dto"
	"orchestrator/facade"
	"orchestrator/pb"
	"orchestrator/pkg/helper"
	"orchestrator/pkg/orchestrator"
	"orchestrator/pkg/pkgerror"
)

// BookingService implement gRPC Server
type BookingService struct{}

func (b BookingService) HandleSyncBooking(ctx context.Context, req *pb.BookingRequest) (resp *pb.BookingSyncResponse, err error) {

	// get request ID from meta-data
	requestID, err := helper.GetRequestID(ctx)
	if err != nil {
		return nil, pkgerror.SetGRPCErrorResp(requestID, err)
	}

	// get facade
	facade := orchestrator.GetInstance().GetSyncFacade(facade.SyncBooking)
	if facade == nil {
		err := eris.Wrap(pkgerror.ErrInternalServerError, "Flow not found")
		return nil, pkgerror.SetGRPCErrorResp(requestID, err)
	}

	// execute facade
	response, err := facade.Run(requestID, req)
	if err != nil {
		return nil, pkgerror.SetGRPCErrorResp(requestID, err)
	}

	// Convert response
	if resp, ok := response.(*pb.BookingSyncResponse); ok {
		return resp, nil
	} else {
		err := eris.Wrap(pkgerror.ErrInternalServerError, "Convert fail")
		return nil, pkgerror.SetGRPCErrorResp(requestID, err)
	}

}

func (b BookingService) HandleAsyncBooking(ctx context.Context, req *pb.BookingRequest) (*pb.BookingASyncResponse, error) {
	// get request ID from meta-data
	requestID, err := helper.GetRequestID(ctx)
	if err != nil {
		return nil, pkgerror.SetGRPCErrorResp(requestID, err)
	}

	// get facade
	flow := orchestrator.GetInstance().GetAsyncFacade(facade.AsyncBooking)
	if flow == nil {
		err := eris.Wrap(pkgerror.ErrInternalServerError, "Flow not found")
		return nil, pkgerror.SetGRPCErrorResp(requestID, err)
	}

	// init dto
	reqMsg := dto.BookingMsgDTO{
		FaultInject:        req.FaultInject,
		ProductID:          req.ProductID,
		AsyncFacadeContext: &orchestrator.AsyncFacadeContext{},
	}

	// execute facade in async
	err = flow.Run(requestID, reqMsg)
	if err != nil {
		return nil, pkgerror.SetGRPCErrorResp(requestID, err)
	}
	// response immediately
	return &pb.BookingASyncResponse{
		RequestID: requestID,
	}, err
}
