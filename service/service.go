package service

import (
	"moonspace/repository/types"
	"moonspace/service/payment"
)

type Service struct {
	Category Category
	Product  Product
	Order    Order
	Payment  payment.Payment
}

func NewService(cli any, cfg types.Config, paymentCfg payment.PaymentCfg) *Service {
	return &Service{
		Category: NewCategoryService(cli, cfg),
		Product:  NewProductService(cli, cfg),
		Order:    NewOrderService(cli, cfg, paymentCfg),
	}
}
