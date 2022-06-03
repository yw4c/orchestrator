package facade

import (
	"orchestrator/node"
	"orchestrator/pkg/orchestrator"
	"orchestrator/topic"
)

const (
	AsyncBooking orchestrator.FacadeName = "AsyncBooking"
	SyncBooking  orchestrator.FacadeName = "SyncBooking"
)

// RegisterAsyncBookingFacade
func RegisterAsyncBookingFacade() {

	// First node: Create order
	createOrderPair := orchestrator.TopicNodePair{
		Topic:     topic.CreateOrder,
		AsyncNode: node.CreateOrderAsyncNode(),
	}
	// Second node: Create payment
	createPaymentPair := orchestrator.TopicNodePair{
		Topic:     topic.CreatePayment,
		AsyncNode: node.CreatePaymentAsyncNode(),
	}
	// Setup Rollback
	rollbackPair := &orchestrator.TopicRollbackNodePair{
		Topic: topic.CancelAsyncBooking,
		Node:  node.CancelBooking(),
	}

	// Wrap nodes in facade
	facade := orchestrator.NewAsyncFacade(topic.CancelAsyncBooking)
	facade.Use(createOrderPair)
	facade.Use(createPaymentPair)

	// Start consume topics of nodes
	facade.Consume()
	// Start consume rollback topic
	facade.ConsumeRollback(rollbackPair)

	// Register facade
	orchestrator.GetInstance().RegisterAsyncFacade(AsyncBooking, facade)

}

// RegisterSyncBookingFacade
func RegisterSyncBookingFacade() {
	// Create facade and register nodes.
	flow := orchestrator.NewSyncFacade()
	flow.Use(node.CreateOrderSyncNode()).
		Use(node.CreatePaymentSyncNode())

	// Start consume rollback Topic
	rollbackPair := &orchestrator.TopicRollbackNodePair{
		Topic: topic.CancelSyncBooking,
		Node:  node.CancelBooking(),
	}
	flow.ConsumeRollback(rollbackPair)

	// register facade
	orchestrator.GetInstance().RegisterSyncFacade(SyncBooking, flow)
}
