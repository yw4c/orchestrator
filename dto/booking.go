package dto

import "orchestrator/pkg/orchestrator"

// Booking DTO of MQ message
type BookingMsgDTO struct {
	FaultInject bool
	ProductID   int64 `json:"product_id"`
	OrderID     int64 `json:"order_id"`
	PaymentID   int64 `json:"payment_id"`
	*orchestrator.AsyncFacadeContext
}
