package events

import "github.com/models"

var Payment payment

// UserCreatedPayload is the data for when a user is created
type PaymentPayload struct {
	ConfirmPaymentByDate *models.ConfirmTransactionPayment
	ConfirmPayment 	*models.ConfirmPaymentIn
}

type payment struct {
	handlers []interface{ Handle(PaymentPayload) }
}

// Register adds an event handler for this event
func (u *payment) Register(handler interface{ Handle(PaymentPayload) }) {
	u.handlers = append(u.handlers, handler)
}

// Trigger sends out an event with the payload
func (u payment) Trigger(payload PaymentPayload) {
	for _, handler := range u.handlers {
		go handler.Handle(payload)
	}
}
