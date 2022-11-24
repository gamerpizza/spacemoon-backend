package payment

import (
	"fmt"
	"moonspace/model"
)

type PaymentType byte

const (
	PaymentTypePayPal PaymentType = 0
	PaymentTypeStripe PaymentType = 1
	PaymentTypeGoogle PaymentType = 2
)

type Payment interface {
	InitPayment(cart model.PaymentRequest, currency string) (any, error)
	CaptureOrder(id string) (any, error)
}

type PaymentService interface {
	GetService(pt PaymentType) (Payment, error)
	Payment
}

type PaymentImpl struct {
	paypalClient Payment
	stripeClient Payment
	googleClient Payment
}

func NewPayment(cfg PaymentCfg) PaymentService {
	p := &PaymentImpl{
		paypalClient: NewPayPalClient(cfg.PayPal),
		stripeClient: NewStripeClient(cfg.Stripe),
		googleClient: NewGoogleClient(),
	}

	return p
}

func (p *PaymentImpl) InitPayment(pr model.PaymentRequest, currency string) (any, error) {
	return p.paypalClient.InitPayment(pr, currency)
}

func (p *PaymentImpl) CaptureOrder(id string) (any, error) {
	return p.paypalClient.CaptureOrder(id)
}

func (p *PaymentImpl) GetService(pt PaymentType) (Payment, error) {
	switch pt {
	case PaymentTypePayPal:
		return checkNil(p.paypalClient, "paypal")
	case PaymentTypeStripe:
		return checkNil(p.stripeClient, "stripe")
	case PaymentTypeGoogle:
		return checkNil(p.googleClient, "google")
	}

	return nil, fmt.Errorf("payment service %d unsupported", pt)
}

func checkNil(ps Payment, pt string) (Payment, error) {
	if ps == nil {
		return nil, fmt.Errorf("%s payment service not initialized", pt)
	}

	return ps, nil
}
