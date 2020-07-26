package handler

import (
	"encoding/json"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog/log"
	"orchestrator/pb"
	"orchestrator/pkg/ctx"
	"orchestrator/pkg/orchestrator"
	"orchestrator/pkg/pkgerror"
)

/***********************************************************
*   異步事務節點
************************************************************/

// Booking 異步流程的推播訊息格式
type ExampleBookingMsg struct {
	orchestrator.AsyncFlowMsg
	OrderID int
	PaymentID int
}

// 建立訂單-異步的事務節點
func CreateOrderAsync() orchestrator.AsyncHandler {
	return func(topic orchestrator.Topic, data []byte) {
		d := &orchestrator.AsyncFlowMsg{}
		if err := json.Unmarshal(data, d); err != nil {
			//log.Println("CreateOrderAsync: ","json Unmarshal fail", err.Error(), string(debug.Stack()))
		}

		// do something


		// has next
		var nextTopic orchestrator.Topic
		if len(d.Topics)-1 >= d.CurrentIndex {
			nextTopic = orchestrator.Topic(d.Topics[d.CurrentIndex+1])
			d.CurrentIndex += 1
			nextData, err := json.Marshal(d)
			if err != nil {
				//log.Println("CreateOrderAsync: ","Marshal fail", string(debug.Stack()))
			}
			orchestrator.GetMQInstance().Produce(nextTopic, nextData)
		}

		//log.Println(
		//	"handling topic:", topic,
		//	"data:", string(data),
		//	"next:", nextTopic,
		//)
	}
}

// 建立付款-異步的事務節點
func CreatePaymentAsync() orchestrator.AsyncHandler {
	return func(topic orchestrator.Topic, data []byte) {
		//log.Println("topic", topic, "data", string(data))
	}
}

/***********************************************************
*   同步事務節點
************************************************************/


// 建立訂單-同步的事務節點
func CreateOrderSync() orchestrator.SyncHandler {
	return func(requestID string, ctx *ctx.Context) error {

		// 從 Context 取出 gRPC Request
		req, isExist := ctx.Get(orchestrator.BookingSyncPbReq)
		if !isExist {
			return eris.Wrap(pkgerror.ErrInternalError,"Can not get request in first handler")
		}
		var request *pb.BookingRequest
		request, ok := req.(*pb.BookingRequest)
		if !ok {
			return eris.Wrap(pkgerror.ErrInternalError,"Convert Request to Protobuf DTO fail")
		}

		// 模擬建立訂單業務邏輯
		var mockOrderID int64 = 1
		log.Info().
			Str("requestID", requestID).
			Int64("productID", request.ProductID).
			Int64("orderID", mockOrderID).
			Msg("Finish Create Order")

		// 將訂單資訊存入 Proto 物件交給下一個
		resp := &pb.BookingSyncResponse{
			RequestID:            requestID,
			OrderID:              mockOrderID,
			PaymentID:            0,
			FaultInject:	request.FaultInject,
		}
		ctx.Set(orchestrator.BookingSyncPbResp, resp)
		return nil
	}

}

// 建立付款單-同步的事務節點
func CreatePaymentSync() orchestrator.SyncHandler {
	return func(requestID string, ctx *ctx.Context) error {

		// 從 Context 取出 gRPC Response DTO
		input, isExist := ctx.Get(orchestrator.BookingSyncPbResp)
		if !isExist {
			return eris.Wrap(pkgerror.ErrInternalError,"Can not get request in first handler")
		}
		var resp *pb.BookingSyncResponse
		resp, ok := input.(*pb.BookingSyncResponse)
		if !ok {
			return eris.Wrap(pkgerror.ErrInternalError,"Convert Request to Protobuf DTO fail")
		}

		// 故障注入
		if resp.FaultInject {
			return eris.Wrap(pkgerror.ErrInvalidInput, "Fault Injection")
		}

		// 模擬建立付款單業務邏輯
		var mockPaymentID int64 = 1
		log.Info().
			Str("requestID", requestID).
			Int64("orderID", resp.OrderID).
			Int64("PaymentID", mockPaymentID).
			Msg("Finish Payment Order")

		resp.PaymentID = mockPaymentID
		return nil
	}
}

/***********************************************************
*  Rollback 事務節點
************************************************************/

func CancelBooking() orchestrator.AsyncHandler {
	return func(topic orchestrator.Topic, data []byte) {
		log.Info().
			Str("topic", string(topic)).
			Str("Message", string(data)).
			Msg("Cancelling Booking, Rollback flow")

	}
}

