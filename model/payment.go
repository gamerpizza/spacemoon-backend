package model

type googlePayCardInfo struct {
	CardDetails string `json:"cardDetails"`
	CardNetwork string `json:"cardNetwork"`
}

type googlePayTokenData struct {
	Token string `json:"token"`
	Type  string `json:"type"`
}

type googlePayPaymentMethodData struct {
	Description      string             `json:"description"`
	Info             googlePayCardInfo  `json:"info"`
	TokenizationData googlePayTokenData `json:"tokenizationData"`
}

type googlePaymentRequest struct {
	Email             string                     `json:"email"`
	PaymentMethodData googlePayPaymentMethodData `json:"paymentMethodData"`
	ShippingAddress   ShippingData               `json:"shippingAddress"`
}

type PaymentRequest struct {
	Order     Order                `json:"order"`
	GooglePay googlePaymentRequest `json:"paymentData"`
}
