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
  - [Concept](#concept)
  - [Task](#task)
    
## Quick Start

````sh
    
    export OCH_PATH=$(pwd)
    cp app.dev.yaml app.yaml
    go test ./... 

    # 本地啟動相關服務
    cd deployment/local
    docker-compose up --build

````
    
## Structure
````
.
├── facade: 註冊同/異步事務流程
├── node: 同/異步事務節點
├── handler: 實現 protobuf gRPC ServiceServer 接口
├── topic: 註冊 Message Queue Topics
````

## Concept
### 節點 (Node)

### 同步流程 (Sync Flow)
![seq-sync](https://www.websequencediagrams.com/cgi-bin/cdraw?lz=dGl0bGUgU2FnYVN5bmMKQ2xpZW50LT5FbmRwb2ludDog6KiC5Zau6LOH5paZIChodHRwKQoAFggtPk9yY2hlc3RyYXRvcgAcEGdSUEMpCgAWDC0-T3JkZQAQGWRlABgK5bu656uLAHcGABEKAFsS57eo6JmfAFcWUGF5bWUAgToKAB4NIAoAFwcAHAsAbgbku5jmrL7llq4AGAoAgVgOABsG6LOH6KiKAIFQFgCCNAoAGBQAgjQKAIJoBgBCEACCXgYKCgoK&s=roundgreen)

* 使用 Context 組合從各服務取得的資料

### 取消流程 (Rollback)
![seq-rollback](https://www.websequencediagrams.com/cgi-bin/cdraw?lz=dGl0bGUgUm9sbGJhY2sKQ2xpZW50LT5FbmRwb2ludDog6KiC5Zau6LOH5paZIChodHRwKQoAFggtPk9yY2hlc3RyYXRvcgAcEGdSUEMpCgAWDC0-T3JkZQAQGWRlABgK5bu656uLAHcGABEKAFsS57eo6JmfAFcWUGF5bWUAgToKAB4NIAoAFwcAHAsAbgbmlLbmrL7llq7lpLHmlZcAHgoAgV4OZXJyb3IAgU8WTVE6IHRyYWNlSUQgKHB1Ymxpc2gpCk1RAIIfEAAbCWNvbnN1bWUAghYXAEIJAIITFY-W5raIAIIMHE9LAIJ1FgCDWQoAgT4NAINSCgCEBgYAgWEJAIN1BgoK&s=roundgreen)
* 發生錯誤時推播 Rollback topic 觸發事件。
* msg 中包含貫穿全路的 traceID (or x-request-id)
* 再以同步方式向各服務發起取消要求，各服務以 traceID 取消資料

### 異步流程 (Async Flow)
![sql-async](https://www.websequencediagrams.com/cgi-bin/cdraw?lz=dGl0bGUgU2FnYUFzeW5jCgpDbGllbnQtPkVuZHBvaW50OiDoqILllq7os4fmlpkgKGh0dHApCgAWCC0-T3JjaGVzdHJhdG9yABwQZ1JQQykKCgAXDC0-TVEAQxBwdWJsaXNoKQAaDwB4Ck9LAEQIAG4KAIEiBgAWBgCBDgYKTVEAcR5jb25zdW1lAFgQT3JkZQCBJBcAFgUAGQnlu7rnq4sAggkGABMIAIB_CENhbGxiYWNrAIIZCAAxCQCCEgwAgT4KCgoKCgoK&s=roundgreen)
* 資料以 json 格式在 MQ 傳遞，orchestrator(replicas) 自行推播接收。再轉 gRPC 向服務請求。
* 接收請求後直接 response , 流程背景處理
* 資料 response 以 callback 方式回覆

### 節流流程 (Throttling Flow)
![seq-throttling](https://www.websequencediagrams.com/cgi-bin/cdraw?lz=dGl0bGUgVGhyb3R0bGluZwpDbGllbnQtPkVuZHBvaW50OiDoqILllq7os4fmlpkgKGh0dHApCgAWCC0-T3JjaGVzdHJhdG9yABwQZ1JQQykKABYMACAQ55uj6IG95rWB56iL5piv5ZCm5a6M5oiQACYPTVEAdxBwdWJsaXNoKQpNUQBmHmNvbnN1bWUAdRJkZQCBFxlkZQAYCuW7uueriwCBfgYAEQoAgWISSUQAgUwk5L-u5pS5AIFqBueCuuW3sgCBWxUAgm0QAFMKAIJvCgCDIwYAcwwAgxUG&s=roundgreen)
* 使用 MQ 透過 Queue 與控制 Concurrency 數，改善同步流程造成 upstream 過載情形，達成 熔斷, retry, rate limit 無法做到的確保執行。
* 將請求連線 lock 住，後續異步執行。
* 使用 topic + pod_id 達成同 pod 進出 MQ
![flow-throttling](https://mermaid.ink/img/eyJjb2RlIjoiZ3JhcGggTFJcbiAgICBHYXRld2F5XG4gICAgTVFcbiAgICBHYXRld2F5IC0tIDEuIGdSUEMgTG9hZEJhbGFuY2luZyAtLT5wb2QxXG4gICAgR2F0ZXdheS0tIDEuIGdSUEMgTG9hZEJhbGFuY2luZyAtLT4gcG9kMlxuICAgIHBvZDEtLSAyLiBwdWIgdG9waWMucG9kMSAtLT5NUVxuICAgIHBvZDItLSAyLiBwdWIgdG9waWMucG9kMiAtLT5NUVxuICAgIE1RIC0tIDMuIHN1YiB0b3BpYy5wb2QxIC0tPiBwb2QxXG4gICAgTVEgLS0gMy4gc3ViIHRvcGljLnBvZDIgLS0-IHBvZDJcbiAgICBcbiAgICBzdWJncmFwaCBvcmNoZXN0cmF0b3JcbiAgICBwb2QxXG4gICAgcG9kMlxuICAgIGVuZFxuIiwibWVybWFpZCI6eyJ0aGVtZSI6ImRlZmF1bHQiLCJ0aGVtZVZhcmlhYmxlcyI6eyJiYWNrZ3JvdW5kIjoid2hpdGUiLCJwcmltYXJ5Q29sb3IiOiIjRUNFQ0ZGIiwic2Vjb25kYXJ5Q29sb3IiOiIjZmZmZmRlIiwidGVydGlhcnlDb2xvciI6ImhzbCg4MCwgMTAwJSwgOTYuMjc0NTA5ODAzOSUpIiwicHJpbWFyeUJvcmRlckNvbG9yIjoiaHNsKDI0MCwgNjAlLCA4Ni4yNzQ1MDk4MDM5JSkiLCJzZWNvbmRhcnlCb3JkZXJDb2xvciI6ImhzbCg2MCwgNjAlLCA4My41Mjk0MTE3NjQ3JSkiLCJ0ZXJ0aWFyeUJvcmRlckNvbG9yIjoiaHNsKDgwLCA2MCUsIDg2LjI3NDUwOTgwMzklKSIsInByaW1hcnlUZXh0Q29sb3IiOiIjMTMxMzAwIiwic2Vjb25kYXJ5VGV4dENvbG9yIjoiIzAwMDAyMSIsInRlcnRpYXJ5VGV4dENvbG9yIjoicmdiKDkuNTAwMDAwMDAwMSwgOS41MDAwMDAwMDAxLCA5LjUwMDAwMDAwMDEpIiwibGluZUNvbG9yIjoiIzMzMzMzMyIsInRleHRDb2xvciI6IiMzMzMiLCJtYWluQmtnIjoiI0VDRUNGRiIsInNlY29uZEJrZyI6IiNmZmZmZGUiLCJib3JkZXIxIjoiIzkzNzBEQiIsImJvcmRlcjIiOiIjYWFhYTMzIiwiYXJyb3doZWFkQ29sb3IiOiIjMzMzMzMzIiwiZm9udEZhbWlseSI6IlwidHJlYnVjaGV0IG1zXCIsIHZlcmRhbmEsIGFyaWFsIiwiZm9udFNpemUiOiIxNnB4IiwibGFiZWxCYWNrZ3JvdW5kIjoiI2U4ZThlOCIsIm5vZGVCa2ciOiIjRUNFQ0ZGIiwibm9kZUJvcmRlciI6IiM5MzcwREIiLCJjbHVzdGVyQmtnIjoiI2ZmZmZkZSIsImNsdXN0ZXJCb3JkZXIiOiIjYWFhYTMzIiwiZGVmYXVsdExpbmtDb2xvciI6IiMzMzMzMzMiLCJ0aXRsZUNvbG9yIjoiIzMzMyIsImVkZ2VMYWJlbEJhY2tncm91bmQiOiIjZThlOGU4IiwiYWN0b3JCb3JkZXIiOiJoc2woMjU5LjYyNjE2ODIyNDMsIDU5Ljc3NjUzNjMxMjglLCA4Ny45MDE5NjA3ODQzJSkiLCJhY3RvckJrZyI6IiNFQ0VDRkYiLCJhY3RvclRleHRDb2xvciI6ImJsYWNrIiwiYWN0b3JMaW5lQ29sb3IiOiJncmV5Iiwic2lnbmFsQ29sb3IiOiIjMzMzIiwic2lnbmFsVGV4dENvbG9yIjoiIzMzMyIsImxhYmVsQm94QmtnQ29sb3IiOiIjRUNFQ0ZGIiwibGFiZWxCb3hCb3JkZXJDb2xvciI6ImhzbCgyNTkuNjI2MTY4MjI0MywgNTkuNzc2NTM2MzEyOCUsIDg3LjkwMTk2MDc4NDMlKSIsImxhYmVsVGV4dENvbG9yIjoiYmxhY2siLCJsb29wVGV4dENvbG9yIjoiYmxhY2siLCJub3RlQm9yZGVyQ29sb3IiOiIjYWFhYTMzIiwibm90ZUJrZ0NvbG9yIjoiI2ZmZjVhZCIsIm5vdGVUZXh0Q29sb3IiOiJibGFjayIsImFjdGl2YXRpb25Cb3JkZXJDb2xvciI6IiM2NjYiLCJhY3RpdmF0aW9uQmtnQ29sb3IiOiIjZjRmNGY0Iiwic2VxdWVuY2VOdW1iZXJDb2xvciI6IndoaXRlIiwic2VjdGlvbkJrZ0NvbG9yIjoicmdiYSgxMDIsIDEwMiwgMjU1LCAwLjQ5KSIsImFsdFNlY3Rpb25Ca2dDb2xvciI6IndoaXRlIiwic2VjdGlvbkJrZ0NvbG9yMiI6IiNmZmY0MDAiLCJ0YXNrQm9yZGVyQ29sb3IiOiIjNTM0ZmJjIiwidGFza0JrZ0NvbG9yIjoiIzhhOTBkZCIsInRhc2tUZXh0TGlnaHRDb2xvciI6IndoaXRlIiwidGFza1RleHRDb2xvciI6IndoaXRlIiwidGFza1RleHREYXJrQ29sb3IiOiJibGFjayIsInRhc2tUZXh0T3V0c2lkZUNvbG9yIjoiYmxhY2siLCJ0YXNrVGV4dENsaWNrYWJsZUNvbG9yIjoiIzAwMzE2MyIsImFjdGl2ZVRhc2tCb3JkZXJDb2xvciI6IiM1MzRmYmMiLCJhY3RpdmVUYXNrQmtnQ29sb3IiOiIjYmZjN2ZmIiwiZ3JpZENvbG9yIjoibGlnaHRncmV5IiwiZG9uZVRhc2tCa2dDb2xvciI6ImxpZ2h0Z3JleSIsImRvbmVUYXNrQm9yZGVyQ29sb3IiOiJncmV5IiwiY3JpdEJvcmRlckNvbG9yIjoiI2ZmODg4OCIsImNyaXRCa2dDb2xvciI6InJlZCIsInRvZGF5TGluZUNvbG9yIjoicmVkIiwibGFiZWxDb2xvciI6ImJsYWNrIiwiZXJyb3JCa2dDb2xvciI6IiM1NTIyMjIiLCJlcnJvclRleHRDb2xvciI6IiM1NTIyMjIiLCJjbGFzc1RleHQiOiIjMTMxMzAwIiwiZmlsbFR5cGUwIjoiI0VDRUNGRiIsImZpbGxUeXBlMSI6IiNmZmZmZGUiLCJmaWxsVHlwZTIiOiJoc2woMzA0LCAxMDAlLCA5Ni4yNzQ1MDk4MDM5JSkiLCJmaWxsVHlwZTMiOiJoc2woMTI0LCAxMDAlLCA5My41Mjk0MTE3NjQ3JSkiLCJmaWxsVHlwZTQiOiJoc2woMTc2LCAxMDAlLCA5Ni4yNzQ1MDk4MDM5JSkiLCJmaWxsVHlwZTUiOiJoc2woLTQsIDEwMCUsIDkzLjUyOTQxMTc2NDclKSIsImZpbGxUeXBlNiI6ImhzbCg4LCAxMDAlLCA5Ni4yNzQ1MDk4MDM5JSkiLCJmaWxsVHlwZTciOiJoc2woMTg4LCAxMDAlLCA5My41Mjk0MTE3NjQ3JSkifX19)


## Task

### 1. 實現事務節點
* 在 ./node 下註冊節點
#### 同步事務節點
同步事務將 gRPC response 物件放進 Context 傳遞資料。我們定義 key 方便取得。
```go
const (
	BookingSyncPbReq orchestrator.FlowContextKeyReq = "bookingSyncPbReq"
	BookingSyncPbResp  = "bookingSyncPbResp"
)
```
取得參數：將 gRPC Request/Response 從 context 取出，轉回 protobuf 物件
```go
func CreateOrderSync() orchestrator.SyncNode {
	return func(requestID string, ctx *ctx.Context) error {
        // 從 Context 取出 gRPC Request
        if req, isExist := ctx.Get(string(BookingSyncPbReq)) ; isExist {
            if request, ok := req.(*pb.BookingRequest); ok {  
```
將處理後的資料放入 gRPC Response 再塞進 Context 繼續傳遞
```go
    resp := &pb.BookingSyncResponse{
        RequestID:            requestID,
        OrderID:              mockOrderID,
        PaymentID:            0,
    }
    ctx.Set(string(BookingSyncPbResp), resp)
```
#### 異步事務節點
定義 flow 傳遞資料的物件，需繼承 ``` *orchestrator.AsyncFlowContext ```
```go
type BookingMsgDTO struct {
	//**** 請求參數 ****//
	FaultInject bool
	ProductID   int64 `json:"product_id"`
	//**** 傳遞資料 ****//
	OrderID int64 `json:"order_id"`
	PaymentID int64 `json:"payment_id"`
    //**** Context ****//
	*orchestrator.AsyncFlowContext
}
```
AsyncNode 會注入以下參數
* topic topic ID
* data marshaled json data
* next 進入下一節點
* rollback 發生錯誤，取消 flow
```go
func CreateOrderAsync() orchestrator.AsyncNode {
	return func(topic orchestrator.Topic, data []byte, next orchestrator.Next, rollback orchestrator.Rollback) {
```
接收參數：將 json decode 取得 Context
```go
    d := &BookingMsgDTO{}
    if err := json.Unmarshal(data, d); err != nil {
```
將處理後的資料放入 context, 調用 next() 送給下一個節點
```go
d.PaymentID = mockPaymentID
next(d)
```
發生錯誤調用 rollback()
````go
    if err := nil {
        rollback(err, d)
        return
    }
````


* Example 請參考 [node](./node/booking.go)

### 2. 註冊事務流程
* 在 ./facade 下註冊流程
* 宣告 facade，方便後續調用 flow
```go
const (
	AsyncBooking orchestrator.Facade = "AsyncBooking"
	SyncBooking  = "SyncBooking"
	ThrottlingBooking  = "ThrottlingBooking"
)
```

建立 Pairs ：一個 Topic 對應一個節點, 同步事務只會有一個 pair : rollback
```go
	createPaymentPair := orchestrator.TopicNodePair{
		Topic:     topic.CreatePayment,
		AsyncNode: node.CreatePaymentAsync(),
	}
```
建立 Flow
```
// 異步
flow := orchestrator.NewAsyncFlow(topic.CancelAsyncBooking)
// 同步
flow := orchestrator.NewSyncFlow()
// 節流
flow := orchestrator.NewThrottlingFlow(topic.CancelThrottlingBooking)
```

註冊 Pairs，將會依序執行節點
```go
	flow.Use(createOrderPair)
	flow.Use(createPaymentPair)
```

監聽 Topic: 啟用所有 consumers
```go
	// 開始監聽異步事務 Topic
	flow.Consume()
	// 開始監聽 rollback topic
	flow.ConsumeRollback(rollbackPair)
```

緩存流程，供 handler 協程調用
```go
orchestrator.GetInstance().SetThrottlingFlows(ThrottlingBooking, flow)
```


* Example 請參考 [facade](./facade/booking.go)

### 3. 調用事務流程
在協程中，從緩存取出註冊過的 flow, 調用 Run() 開始執行
````go
flow := orchestrator.GetInstance().GetThrottlingFlows(facade.ThrottlingBooking)
flow.Run(requestID, reqMsg)
````

* Example 請參考 [handler](./handler/booking.go)