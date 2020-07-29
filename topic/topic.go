package topic

import (
	"orchestrator/pkg/orchestrator"
)

const (
	CreateOrder        orchestrator.Topic = "CreateOrder"
	CreatePayment      orchestrator.Topic = "CreatePayment"
	CancelAsyncBooking orchestrator.Topic = "CancelAsyncBooking"
	CancelSyncBooking orchestrator.Topic = "CancelSyncBooking"
)
