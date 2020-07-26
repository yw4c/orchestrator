package topic

import (
	"orchestrator/pkg/orchestrator"
)

const (
	CreateOrder   orchestrator.Topic = "CreateOrder"
	CreatePayment orchestrator.Topic = "CreatePayment"
	CancelBooking orchestrator.Topic = "CancelBooking"
)
