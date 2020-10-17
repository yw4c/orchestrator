# Orchestrator Framework of Saga pattern
* Orchestrator 在微服務中，負責集中處理分散事務流程; 上游服務只需維護 gRPC 接口。
* 流程支援同步(Sync)/異步(Async)，以及鎖住連線使用隊列的節流(Throttling)模式。
* 只要實現框架接口，可以輕易註冊、擴展分散事務流程。
* 可抽換 Message Queue ，目前支援 RabbitMQ, NATS。


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
本地啟動相關服務
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

## Concept
### 同步流程 (Sync Flow)
![Alt text](https://www.websequencediagrams.com/cgi-bin/cdraw?lz=dGl0bGUgU2FnYVN5bmMKQ2xpZW50LT5FbmRwb2ludDog6KiC5Zau6LOH5paZIChodHRwKQoAFggtPk9yY2hlc3RyYXRvcgAcEGdSUEMpCgAWDC0-T3JkZQAQGWRlABgK5bu656uLAHcGABEKAFsS57eo6JmfAFcWUGF5bWUAgToKAB4NIAoAFwcAHAsAbgbku5jmrL7llq4AGAoAgVgOABsG6LOH6KiKAIFQFgCCNAoAGBQAgjQKAIJoBgBCEACCXgYKCgoK&s=roundgreen)

* 使用 Context 組合從各服務取得的資料

### 取消流程 (Rollback)
![Alt text](https://www.websequencediagrams.com/cgi-bin/cdraw?lz=dGl0bGUgUm9sbGJhY2sKQ2xpZW50LT5FbmRwb2ludDog6KiC5Zau6LOH5paZIChodHRwKQoAFggtPk9yY2hlc3RyYXRvcgAcEGdSUEMpCgAWDC0-T3JkZQAQGWRlABgK5bu656uLAHcGABEKAFsS57eo6JmfAFcWUGF5bWUAgToKAB4NIAoAFwcAHAsAbgbmlLbmrL7llq7lpLHmlZcAHgoAgV4OZXJyb3IAgU8WTVE6IHRyYWNlSUQgKHB1Ymxpc2gpCk1RAIIfEAAbCWNvbnN1bWUAghYXAEIJAIITFY-W5raIAIIMHE9LAIJ1FgCDWQoAgT4NAINSCgCEBgYAgWEJAIN1BgoK&s=roundgreen)
* 為支援同、異步；發生錯誤時推播 Rollback topic 觸發事件。
* msg 中包含貫穿全路的 traceID (or x-request-id)
* 再以同步方式向各服務發起取消要求，各服務以 traceID 取消資料

### 異步流程 (Async Flow)
![Alt text](https://www.websequencediagrams.com/cgi-bin/cdraw?lz=dGl0bGUgU2FnYUFzeW5jCgpDbGllbnQtPkVuZHBvaW50OiDoqILllq7os4fmlpkgKGh0dHApCgAWCC0-T3JjaGVzdHJhdG9yABwQZ1JQQykKCgAXDC0-TVEAQxBwdWJsaXNoKQAaDwB4Ck9LAEQIAG4KAIEiBgAWBgCBDgYKTVEAcR5jb25zdW1lAFgQT3JkZQCBJBcAFgUAGQnlu7rnq4sAggkGABMIAIB_CENhbGxiYWNrAIIZCAAxCQCCEgwAgT4KCgoKCgoK&s=roundgreen)
* 資料以 json 格式在 MQ 傳遞，orchestrator(replicas) 自行推播接收。再轉 gRPC 向服務請求。
* 接收請求後直接 response , 流程背景處理
* 資料 response 以 callback 方式回覆

### 節流流程 (Throttling Flow)
![Alt text](https://www.websequencediagrams.com/cgi-bin/cdraw?lz=dGl0bGUgU2FnYUh5YnJpZApDbGllbnQtPkVuZHBvaW50OiDoqILllq7os4fmlpkgKGh0dHApCgAWCC0-T3JjaGVzdHJhdG9yABwQZ1JQQykKABYMLT5SZWRpczog55uj6IG95rWB56iL5piv5ZCm5a6M5oiQAB8PTVEAcBBwdWJsaXNoKQpNUQBfHmNvbnN1bWUAcBBPcmRlAIEQGWRlABgK5bu656uLAIF3BgARCgCBYQxPSwCBRh3kv67mlLkAgV0G54K65beyAIFcBwCBfgUAgikRp7jnmbwAgh4PAIJ7CgBiCgCCcQoAgyUGAIECBgCDEQY&s=roundgreen)
* 使用 MQ 透過 Queue 與控制 Concurrency 數，改善同步流程造成 upstream 過載情形，達成 熔斷, retry, rate limit 無法做到的確保執行。
* 將請求連線 lock 住，後續異步執行。
* orchestrator 有多個 replica , 所以我們須 watch redis 流程是否已完成。


## Task


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
