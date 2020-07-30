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
type BookingMsgDTO struct {

	//**** 請求參數 ****//
	FaultInject bool
	ProductID   int64 `json:"product_id"`

	//**** 傳遞資料 ****//
	OrderID int64 `json:"order_id"`
	PaymentID int64 `json:"payment_id"`

	*orchestrator.AsyncFlowContext
}

// 建立訂單-異步的事務節點
func CreateOrderAsync() orchestrator.AsyncHandler {
	return func(topic orchestrator.Topic, data []byte, next orchestrator.Next, rollback orchestrator.Rollback) {

		// Convert Message into struct
		d := &BookingMsgDTO{}

		if err := json.Unmarshal(data, d); err != nil {
			rollback(eris.Wrap(pkgerror.ErrInternalError, "json unmarshal fail"), d)
			return
		}
		// 模擬建立訂單業務邏輯
		var mockOrderID int64 = 11
		d.OrderID = mockOrderID

		next(d)

		log.Info().Msg("CreateOrderAsync finished")
	}
}

// 建立付款-異步的事務節點
func CreatePaymentAsync() orchestrator.AsyncHandler {
	return func(topic orchestrator.Topic, data []byte, next orchestrator.Next, rollback orchestrator.Rollback) {

		// Convert Message into struct
		d := &BookingMsgDTO{}

		if err := json.Unmarshal(data, d); err != nil {
			rollback(eris.Wrap(pkgerror.ErrInternalError, "json unmarshal fail"), d)
			return
		}

		// 模擬建立訂單業務邏輯
		var mockPaymentID int64 = 12
		d.PaymentID = mockPaymentID

		// 故障注入
		if d.FaultInject {
			rollback(eris.Wrap(pkgerror.ErrInvalidInput, "this is an mocked invalid error"), d)
			return
		}

		next(d)

		log.Info().Msg("CreatePaymentAsync finished")
	}
}

/***********************************************************
*   同步事務節點
************************************************************/

// 事務節點間，傳遞 Context 資料的 key
const (
	BookingSyncPbReq orchestrator.FlowContextKeyReq = "bookingSyncPbReq"
	BookingSyncPbResp orchestrator.FlowContextKeyResp = "bookingSyncPbResp"
)

// 建立訂單-同步的事務節點
func CreateOrderSync() orchestrator.SyncHandler {
	return func(requestID string, ctx *ctx.Context) error {

		// 從 Context 取出 gRPC Request
		req, isExist := ctx.Get(string(BookingSyncPbReq))
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
		ctx.Set(string(BookingSyncPbResp), resp)
		return nil
	}

}

// 建立付款單-同步的事務節點
func CreatePaymentSync() orchestrator.SyncHandler {
	return func(requestID string, ctx *ctx.Context) error {

		// 從 Context 取出 gRPC Response DTO
		input, isExist := ctx.Get(string(BookingSyncPbResp))
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

func CancelBooking() orchestrator.RollbackHandler {
	return func(topic orchestrator.Topic, data []byte) {

		msg := &orchestrator.RollbackMsg{}
		if err := json.Unmarshal(data, msg); err != nil {
			log.Panic().Str("topic", string(topic)).Str("data", string(data)).Msg("Failed to unmarshal")
			return
		}

		// 模擬 Cancel
		// go orderClient.Cancel(msg.RequestID)
		// go paymentClient.Cancel(msg.RequestID)

		log.Info().
			Str("topic", string(topic)).
			Str("Message", string(data)).
			Str("requestID", msg.RequestID).
			Msg("Cancelling Booking, Rollback flow")

	}
}

