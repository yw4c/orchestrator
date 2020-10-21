package facade

import (
	"orchestrator/node"
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
	createOrderPair := orchestrator.TopicNodePair{
		Topic:     topic.CreateOrder,
		AsyncNode: node.CreateOrderAsync(),
	}
	// 建立付款單
	createPaymentPair := orchestrator.TopicNodePair{
		Topic:     topic.CreatePayment,
		AsyncNode: node.CreatePaymentAsync(),
	}
	// Rollback
	rollbackPair := &orchestrator.TopicRollbackNodePair{
		Topic: topic.CancelAsyncBooking,
		Node:  node.CancelBooking(),
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
	flow.Use(node.CreateOrderSync()).
		Use(node.CreatePaymentSync())

	// 開始監聽 rollback Topic
	rollbackPair := &orchestrator.TopicRollbackNodePair{
		Topic: topic.CancelSyncBooking,
		Node:  node.CancelBooking(),
	}
	flow.ConsumeRollback(rollbackPair)

	// 註冊流程
	orchestrator.GetInstance().SetSyncFlows(SyncBooking, flow)
}

func RegisterThrottlingBookingFlow()  {

	// 建立訂單
	createOrderPair := orchestrator.TopicNodePair{
		Topic:     topic.CreateOrderThrottling,
		AsyncNode: node.CreateOrderAsync(),
	}
	// 建立付款單
	createPaymentPair := orchestrator.TopicNodePair{
		Topic:     topic.CreatePaymentThrottling,
		AsyncNode: node.CreatePaymentAsync(),
	}
	// Rollback
	rollbackPair := &orchestrator.TopicRollbackNodePair{
		Topic: topic.CancelAsyncBooking,
		Node:  node.CancelBooking(),
	}

	flow := orchestrator.NewThrottlingFlow(topic.CancelThrottlingBooking)
	flow.Use(createOrderPair)
	flow.Use(createPaymentPair)
	flow.ConsumeRollback(rollbackPair)
	flow.Consume()

	orchestrator.GetInstance().SetThrottlingFlows(ThrottlingBooking, flow)

}