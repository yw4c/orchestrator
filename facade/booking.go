package facade

import (
	"orchestrator/handler"
	"orchestrator/pkg/orchestrator"
	"orchestrator/topic"
)

const (
	AsyncBooking orchestrator.Facade = "AsyncBooking"
	SyncBooking orchestrator.Facade = "SyncBooking"
	ThrottlingBooking orchestrator.Facade = "ThrottlingBooking"
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

func RegisterThrottlingBookingFlow()  {

	// 建立訂單
	createOrderPair := orchestrator.TopicHandlerPair{
		Topic:        topic.CreateOrderThrottling,
		AsyncHandler: handler.CreateOrderAsync(),
	}
	// 建立付款單
	createPaymentPair := orchestrator.TopicHandlerPair{
		Topic:        topic.CreatePaymentThrottling,
		AsyncHandler: handler.CreatePaymentAsync(),
	}
	// Rollback
	rollbackPair := &orchestrator.TopicRollbackHandlerPair{
		Topic:        topic.CancelAsyncBooking,
		Handler: handler.CancelBooking(),
	}

	flow := orchestrator.NewThrottlingFlow(topic.CancelThrottlingBooking)
	flow.Use(createOrderPair)
	flow.Use(createPaymentPair)
	flow.ConsumeRollback(rollbackPair)
	flow.Consume()

	orchestrator.GetInstance().SetThrottlingFlows(ThrottlingBooking, flow)

}