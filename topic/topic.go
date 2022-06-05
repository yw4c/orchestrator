package topic

import (
	"orchestrator/pkg/orchestrator"
)

const (

	// Sync
	CancelSyncBooking orchestrator.Topic = "cancel_sync_booking"

	// Async
	CreateOrder        = "create_order"
	CreatePayment      = "create_payment"
	CancelAsyncBooking = "cancel_async_booking"

)
