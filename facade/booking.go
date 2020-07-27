package facade

import (
	"orchestrator/handler"
	"orchestrator/pkg/orchestrator"
	"orchestrator/topic"
)

const (
	AsyncBooking orchestrator.Facade = "AsyncBooking"
	SyncBooking orchestrator.Facade = "SyncBooking"
)


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
	// 建立流程
	flow := orchestrator.NewAsyncFlow(topic.CancelBooking)
	flow.Use(createOrderPair)
	flow.Use(createPaymentPair)

	// 開始監聽異步事務 Topic
	flow.Consume()

	// 註冊流程
	orchestrator.GetInstance().SetAsyncFlows(AsyncBooking, flow)

}

// 註冊同步的 Booking 事務流程
func RegisterSyncBookingFlow() {
	// 建立流程
	flow := orchestrator.NewSyncFlow()
	flow.Use(handler.CreateOrderSync()).
		Use(handler.CreatePaymentSync())

	// 開始監聽 rollback Topic
	rollbackPair := &orchestrator.TopicRollbackHandlerPair{
		Topic:        topic.CancelBooking,
		Handler: handler.CancelBooking(),
	}
	flow.ConsumeRollback(rollbackPair)

	// 註冊流程
	orchestrator.GetInstance().SetSyncFlows(SyncBooking, flow)
}