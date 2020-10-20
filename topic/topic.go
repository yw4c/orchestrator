package topic

import (
	"orchestrator/pkg/orchestrator"
)


const (
	CreateOrder        orchestrator.Topic = "create_order"
	CreatePayment       = "create_payment"
	CancelAsyncBooking  = "cancel_async_booking"
	CancelSyncBooking   = "cancel_sync_booking"
	CreateOrderThrottling = "create_order_throttling"
	CreatePaymentThrottling = "create_payment_throttling"
)
