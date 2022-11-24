package payment

import (
	"fmt"
	"moonspace/model"
)

type GooglePayClient struct{}

func NewGoogleClient() Payment {
	return &GooglePayClient{}
}

func (gp *GooglePayClient) InitPayment(pr model.PaymentRequest, currency string) (any, error) {
	fmt.Println("Request received: ", pr)
	return nil, nil
}

func (gp *GooglePayClient) CaptureOrder(id string) (any, error) {
	// NOOP
	return nil, nil
}
