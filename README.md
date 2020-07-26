# Orchestrator Framework of Saga pattern
* Orchestrator 在微服務中，負責集中處理分散事務流程，
* 實現框架接口，可以輕易註冊，擴展分散事務流程。
* 以 gRPC 作為服務入口。
* 可抽換 Message Queue ，目前支援 RabbitMQ


## Contents

- [Orchestrator Framework of Saga pattern](#orchestrator-framework-of-saga-pattern)
  - [Contents](#contents)
  - [Installation](#installation)
  - [Structure](#structure)
  - [Examples](#examples)
    - [範例 Booking 流程](#範例-booking-流程)
    - [同步事務流程](#同步事務流程)
    - [異步事務流程](#異步事務流程)
    
    
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
* 我們會在請求參數做故障注入，觸發 Rollback 流程

![Alt text](https://www.websequencediagrams.com/cgi-bin/cdraw?lz=dGl0bGUgTGFiIC0gU2FnYSBQYXR0ZXJuCgpDbGllbnQtPkVuZHBvaW50OiBvcmRlckluZm8oaHR0cCkKbm90ZSBsZWZ0IG9mIAAdCkVuZHBpbnQg5Y-v6IO95o6l5pS25L6G6IeqU2VydmVyIFxuIOaIluaYr0Zyb250ZW5kIOeahOiri-axgiAKAGcILT5PcmNoZXN0cmF0b3I6IGdSUEMKCmFsdCBzeW5jCiAgICAAFQwtPk9yZGVyAIEaCywgdHJhY2VJRCAoZ1JQQykALAdkZQAiCuW7uueri-ioguWWrgARDgBuDACBcQYAOA8AbgxQYXltZQCCFwpEAGUVAB4HACMLAHgG5pS25qy-5Zau5aSx5pWXAB4OAIF1DmVycm9y77yIZ1JQQ--8iQCBcRNNUTogcHVibGlzaCB0b3BpY1tyb2xsYmFjazoAggYHXQCCJxMAg1gKAF0FAIIkDACDBAoAhAkGABoJAIN8BmVsc2UgYQCCehcAhC8KT0sANh5PSwBJCACBNR5cbgCBTwdjcmVhdGUtAIUMBToAg2ASAIFfBk1RAIQ0EGNvbnN1bWUAGS8AhCs8AIQbLQCBXA8AgV4OcACEXgcAhFQQAIFJIwAdLACFLC4gAIUOTSAAhTMpAIQACQCFOyQAhTIIQ2EAhgUGIAplbmQAhAMfAIYTKgCIcAcAg1wpawCIJSMAgkgdAIIlFG5vdCBmb3VuZO-8iOacqgCIUweIkOWKn--8iQCIJA8KCgoK&s=napkin)



### 同步事務流程

#### 1. 撰寫事務節點(handler)
* 第一個節點從 Context 中取得 gRPC request DTO 
* 後面以 gRPC response DTO 在 context 中遞送
    
    
```
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

```


* 撰寫 rollback handler, 這是異步節點
```
func CancelBooking() orchestrator.AsyncHandler {
    return func(topic orchestrator.Topic, data []byte) {
        log.Info().
            Str("topic", string(topic)).
            Str("Message", string(data)).
            Msg("Cancelling Booking, Rollback flow")

    }
}
```

* 請參考 ```./hancler/booking.go```
    

#### 1. 註冊事務流程
* 在程式 start up 時進行註冊
* 流程需有一異步的 rollback handler
    
```
    // 註冊同步的 Booking 事務流程
   func RegisterSyncBookingFlow() {
    // 建立流程，依序加入節點
    flow := orchestrator.NewSyncFlow()
    flow.Use(handler.CreateOrderSync()).
        Use(handler.CreatePaymentSync())
   
    // 開始監聽 rollback Topic
    rollbackPair := &orchestrator.TopicHandlerPair{
        Topic:        topic.CancelBooking,
        AsyncHandler: handler.CancelBooking(),
    }
    flow.ConsumeRollback(rollbackPair)
   
    // 註冊流程
    orchestrator.GetInstance().SetFlows(SyncBooking, flow)
   }
```
    
* 請參考 ```./facade/booking.go```

#### 1. 調用事務流程
    * 提供 helper.GetRequestID() 從 metadata 取得 request ID
    * 在請求協程中調用 Run()
    
    ```
	// 從 metadata 取得 request-id
	requestID, err := helper.GetRequestID(ctx)
	if err != nil {
		return nil, pkgerror.SetGRPCErrorResp(requestID, err)
	}

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
    ```

#### 1. 測試事務流程
* Happy Ending
````
grpcurl -rpc-header x-request-id:example-request-id -plaintext -d '{"ProductID": "1", "FaultInject": "false"}' localhost:10000 pb.BookingService/HandleSyncBooking
````
* 模擬建立 付款單 失敗
````
grpcurl -rpc-header x-request-id:example-request-id -plaintext -d '{"ProductID": "1", "FaultInject": "true"}' localhost:10000 pb.BookingService/HandleSyncBooking
````


### 異步事務流程