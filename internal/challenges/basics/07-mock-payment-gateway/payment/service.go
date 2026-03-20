package payment

import (
	"errors"
)

// Gateway is a consumer-side interface for the payment gateway.
// Requirement: Mock this using GoMock or Mockery in your tests.
type Gateway interface {
	Charge(amount int) (string, error) // returns transaction ID
	Refund(id string) error
}

// PaymentService coordinates payments and prevents double-refunds.
type PaymentService struct {
	gateway Gateway
	charged map[string]bool
}

// NewPaymentService creates a new payment service.
func NewPaymentService(g Gateway) *PaymentService {
	return &PaymentService{
		gateway: g,
		charged: make(map[string]bool),
	}
}

// Charge charges the gateway and records the transaction ID if successful.
func (s *PaymentService) Charge(amount int) (string, error) {
	// TODO 1: Call s.gateway.Charge(amount)
	// TODO 2: If successful, record the transaction ID in s.charged
	// TODO 3: Return the ID and any error
	panic("not implemented")
}

// Refund refunds a previously charged transaction.
func (s *PaymentService) Refund(id string) error {
	// TODO 1: Check if 'id' exists in s.charged.
	// TODO 2: If NOT charged, return an error (prevent double-refund or invalid refund).
	// TODO 3: Call s.gateway.Refund(id).
	panic("not implemented")
}

var (
	ErrNotCharged = errors.New("transaction was not charged")
)
