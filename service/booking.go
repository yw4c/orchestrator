package service

import (
	"context"
	"fmt"
	"github.com/rotisserie/eris"
	"orchestrator/facade"
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

	requestID, err := helper.GetRequestID(ctx)
	if err != nil {
		return nil, pkgerror.SetGRPCErrorResp(requestID, err)
	}
	fmt.Println(requestID)

	// 從 facade 取得註冊的事務流程
	flow := orchestrator.GetInstance().GetFlow(facade.SyncBooking)
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

func (b BookingService) HandleAsyncBooking(context.Context, *pb.BookingRequest) (*pb.BookingASyncResponse, error) {
	panic("implement me")
}
