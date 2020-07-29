package service

import (
	"context"
	"github.com/rotisserie/eris"
	"orchestrator/facade"
	"orchestrator/handler"
	"orchestrator/pb"
	"orchestrator/pkg/helper"
	"orchestrator/pkg/orchestrator"
	"orchestrator/pkg/pkgerror"
)

// 協調者對外的 gRPC 服務
// BookingService implement gRPC Server of Protobuf
type BookingService struct {

}

func (b BookingService) HandleSyncBooking(ctx context.Context,req *pb.BookingRequest) (resp *pb.BookingSyncResponse, err error) {

	// 從 metadata 取得 request-id
	requestID, err := helper.GetRequestID(ctx)
	if err != nil {
		return nil, pkgerror.SetGRPCErrorResp(requestID, err)
	}

	// 從 facade 取得註冊的事務流程
	flow := orchestrator.GetInstance().GetSyncFlow(facade.SyncBooking)
	if flow == nil {
		err := eris.Wrap(pkgerror.ErrInternalServerError, "Flow not found")
		return nil, pkgerror.SetGRPCErrorResp(requestID, err)
	}

	// 執行事務流程
	response, err := flow.Run(requestID, req)
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

func (b BookingService) HandleAsyncBooking(ctx context.Context,req *pb.BookingRequest) (*pb.BookingASyncResponse, error) {
	// 從 metadata 取得 request-id
	requestID, err := helper.GetRequestID(ctx)
	if err != nil {
		return nil, pkgerror.SetGRPCErrorResp(requestID, err)
	}

	// 從 facade 取得註冊的事務流程
	flow := orchestrator.GetInstance().GetAsyncFlow(facade.AsyncBooking)
	if flow == nil {
		err := eris.Wrap(pkgerror.ErrInternalServerError, "Flow not found")
		return nil, pkgerror.SetGRPCErrorResp(requestID, err)
	}

	reqMsg := handler.BookingMsgDTO{
		FaultInject:      req.FaultInject,
		ProductID:        req.ProductID,
		AsyncFlowContext: &orchestrator.AsyncFlowContext{},
	}


	// 執行事務流程
	err = flow.Run(requestID, reqMsg)
	if err != nil {
		return nil, pkgerror.SetGRPCErrorResp(requestID, err)
	}
	return &pb.BookingASyncResponse{
		RequestID:            requestID,
		FaultInject:          false,
	}, err
}
