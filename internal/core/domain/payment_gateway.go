package domain

import "context"

// PaymentGateway abstracts external payment providers like Zarinpal
type PaymentGateway interface {
	// RequestPayment initiates a payment request and returns the authority and redirect payment URL
	RequestPayment(ctx context.Context, amountToman int64, description, callbackURL string) (authority string, paymentURL string, err error)
	// VerifyPayment verifies the payment with the gateway and returns the reference ID
	VerifyPayment(ctx context.Context, authority string, amountToman int64) (refID string, err error)
}
