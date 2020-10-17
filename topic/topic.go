package topic

import (
	"orchestrator/pkg/orchestrator"
)


const (
	CreateOrder        orchestrator.Topic = "create_order"
	CreatePayment      orchestrator.Topic = "create_payment"
	CancelAsyncBooking orchestrator.Topic = "cancel_async_booking"
	CancelSyncBooking  orchestrator.Topic = "cancel_sync_booking"
)
