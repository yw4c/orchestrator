# Orchestrator Framework of Saga pattern
* Orchestrator 在微服務中，負責集中處理分散事務流程。
* 流程支援同步/異步。
* 實現框架接口，可以輕易註冊、擴展分散事務流程。
* 以 gRPC 作為服務入口。
* 可抽換 Message Queue ，目前支援 RabbitMQ


## Contents

- [Orchestrator Framework of Saga pattern](#orchestrator-framework-of-saga-pattern)
  - [Contents](#contents)
  - [Quick Start](#quick-start)
  - [Structure](#structure)
  - [Examples](#examples)
    - [範例 Booking 流程](#範例-booking-流程)
    - [同步事務流程](#同步事務流程)
        - [實現事務節點](#實現事務節點)
        - [註冊事務流程](#註冊事務流程)
        - [調用事務流程](#調用事務流程)
        - [測試事務流程](#測試事務流程)
    - [異步事務流程](#異步事務流程)
        - [定義訊息物件](#定義訊息物件)
        - [實現事務節點](#實現事務節點)
        - [註冊事務流程](#註冊事務流程)
        - [調用事務流程](#調用事務流程)
        - [測試事務流程](#測試事務流程)
    - [實現 rollback 事件](#實現-rollback-事件)
    
## Quick Start
本地啟動 RabbitMQ
````
    cd deployment/local
    docker-compose build
    docker-compose up
````
    
## Structure
````
.
├── facade: 註冊同/異步事務流程
├── handler: 同/異步事務節點(handler)
├── service: 實現 protobuf gRPC ServiceServer 接口
├── topic: 註冊 Message Queue Topics
````

## Examples

### 範例 Booking 流程
* 模擬兩服務分別為 訂單Order 和 收款Payment
* 流程為：建立訂單後建立收款單，收款單建立失敗訂單也須註銷
* 範例分別實現同步/異步版本
* 範例使用 x-request-id 做為 traceID，可從 metadata 取得
* 可在請求參數做故障注入，觸發 Rollback 流程

![Alt text](https://www.websequencediagrams.com/cgi-bin/cdraw?lz=dGl0bGUgTGFiIC0gU2FnYSBQYXR0ZXJuCgpDbGllbnQtPkVuZHBvaW50OiBvcmRlckluZm8oaHR0cCkKbm90ZSBsZWZ0IG9mIAAdCkVuZHBpbnQg5Y-v6IO95o6l5pS25L6G6IeqU2VydmVyIFxuIOaIluaYr0Zyb250ZW5kIOeahOiri-axgiAKAGcILT5PcmNoZXN0cmF0b3I6IGdSUEMKCmFsdCBzeW5jCiAgICAAFQwtPk9yZGVyAIEaCywgdHJhY2VJRCAoZ1JQQykALAdkZQAiCuW7uueri-ioguWWrgARDgBuDACBcQYAOA8AbgxQYXltZQCCFwpEAGUVAB4HACMLAHgG5pS25qy-5Zau5aSx5pWXAB4OAIF1DmVycm9y77yIZ1JQQ--8iQCBcRNNUTogcHVibGlzaCB0b3BpY1tyb2xsYmFjazoAggYHXQCCJxMAg1gKAF0FAIIkDACDBAoAhAkGABoJAIN8BmVsc2UgYQCCehcAhC8KT0sANh5PSwBJCACBNR5cbgCBTwdjcmVhdGUtAIUMBToAg2ASAIFfBk1RAIQ0EGNvbnN1bWUAGS8AhCs8AIQbLQCBXA8AgV4OcACEXgcAhFQQAIFJIwAdLACFLC4gAIUOTSAAhTMpAIQACQCFOyQAhTIIQ2EAhgUGIAplbmQAhAMfAIYTKgCIcAcAg1wpawCIJSMAgkgdAIIlFG5vdCBmb3VuZO-8iOacqgCIUweIkOWKn--8iQCIJA8KCgoK&s=napkin)



### 同步事務流程

#### 實現事務節點
* 第一個節點從 Context 中取得 gRPC request DTO 
* 後面以 gRPC response DTO 在 context 中遞送
    
    
```go
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

```
* 請參考 [handler](./handler/booking.go)


    

#### 註冊事務流程
* 在程式 start up 時進行註冊
* 流程需有一異步的 rollback handler
    
```go
// 註冊同步的 Booking 事務流程
func RegisterSyncBookingFlow() {
	// 建立流程
	flow := orchestrator.NewSyncFlow()
	flow.Use(handler.CreateOrderSync()).
		Use(handler.CreatePaymentSync())

	// 開始監聽 rollback Topic
	rollbackPair := &orchestrator.TopicRollbackHandlerPair{
		Topic:        topic.CancelSyncBooking,
		Handler: handler.CancelBooking(),
	}
	flow.ConsumeRollback(rollbackPair)

	// 註冊流程
	orchestrator.GetInstance().SetSyncFlows(SyncBooking, flow)
}
```
    

* 請參考 [facade](./facade/booking.go)

#### 調用事務流程
* 提供 helper.GetRequestID() 從 metadata 取得 request ID
* 在請求協程中調用 Run() 開始執行事務
    
```go
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
	response, err := flow.Run(requestID, req, handler.BookingSyncPbReq, handler.BookingSyncPbResp)
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
```

* 請參考 [service](./service/booking.go)

#### 測試事務流程
* Happy Ending
````
grpcurl -rpc-header x-request-id:example-request-id -plaintext -d '{"ProductID": "1", "FaultInject": "false"}' localhost:10000 pb.BookingService/HandleSyncBooking
````
* 模擬建立 付款單 失敗
````
grpcurl -rpc-header x-request-id:example-request-id -plaintext -d '{"ProductID": "1", "FaultInject": "true"}' localhost:10000 pb.BookingService/HandleSyncBooking
````


### 異步事務流程

#### 定義訊息物件
* 需繼承 ```*orchestrator.AsyncFlowContext```

```go
// Booking 異步流程的推播訊息格式
type BookingMsgDTO struct {

	//**** 請求參數 ****//
	FaultInject bool
	ProductID   int64 `json:"product_id"`

	//**** 傳遞資料 ****//
	OrderID int64 `json:"order_id"`
	PaymentID int64 `json:"payment_id"`

    //**** !!! 必需繼承 Context  !!! ****//
	*orchestrator.AsyncFlowContext
}
```

#### 實現事務節點

* ```next()``` 觸發執行下一節點
* ```rollback()``` 紀錄錯誤日誌，觸發 rollback Topic


```go
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
```

* 請參考 [handler](./handler/booking.go)

    

#### 註冊事務流程
* 在程式 start up 時進行註冊
* 流程需有一異步的 rollback handler
    
```go
// 註冊異步的 Booking 事務流程
func RegisterAsyncBookingFlows() {

	// 建立訂單
	createOrderPair := orchestrator.TopicHandlerPair{
		Topic:        topic.CreateOrder,
		AsyncHandler: handler.CreateOrderAsync(),
	}
	// 建立付款單
	createPaymentPair := orchestrator.TopicHandlerPair{
		Topic:        topic.CreatePayment,
		AsyncHandler: handler.CreatePaymentAsync(),
	}
	// Rollback
	rollbackPair := &orchestrator.TopicRollbackHandlerPair{
		Topic:        topic.CancelAsyncBooking,
		Handler: handler.CancelBooking(),
	}

	// 建立流程
	flow := orchestrator.NewAsyncFlow(topic.CancelAsyncBooking)
	flow.Use(createOrderPair)
	flow.Use(createPaymentPair)

	// 開始監聽異步事務 Topic
	flow.Consume()
	// 開始監聽 rollback topic
	flow.ConsumeRollback(rollbackPair)

	// 註冊流程
	orchestrator.GetInstance().SetAsyncFlows(AsyncBooking, flow)

}
```
    

* 請參考 [facade](./facade/booking.go)

#### 調用事務流程
* HandleAsyncBooking 為實現 gRPC service 的 endpoint
* 提供 helper.GetRequestID() 從 metadata 取得 request ID
* 在請求協程中調用 Run() 開始執行事務
    
```go
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

	// 加入請求
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

```

* 請參考 [service](./service/booking.go)

#### 測試事務流程
* Happy Ending
````
grpcurl -rpc-header x-request-id:example-request-id -plaintext -d '{"ProductID": "1", "FaultInject": "false"}' localhost:10000 pb.BookingService/HandleAsyncBooking
````
* 模擬建立 付款單 失敗
````
grpcurl -rpc-header x-request-id:example-request-id -plaintext -d '{"ProductID": "1", "FaultInject": "true"}' localhost:10000 pb.BookingService/HandleAsyncBooking
````

### 實現 rollback 事件

* 以異步節點方式註冊一 topic
```go
CancelBooking orchestrator.Topic = "CancelBooking"
```

* 請參考 [topic](./topic/booking.go)

* 撰寫 rollback handler
```go
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

```
* 請參考 [handler](./handler/booking.go)
