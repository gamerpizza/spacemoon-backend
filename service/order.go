package service

import (
	"context"
	"moonspace/model"
	"moonspace/repository"
	repo_types "moonspace/repository/types"
	"moonspace/service/payment"
	"time"

	"github.com/stripe/stripe-go/v73"
)

type Order interface {
	CreateOrder(o model.PaymentRequest, paymentType payment.PaymentType) (any, error)
	Checkout(oid string, userId model.UserID, paymentType payment.PaymentType) error
}

type orderImpl struct {
	orderRepo repo_types.Repository[model.Order]
	payment   payment.PaymentService
}

func NewOrderService(cli any, cfg repo_types.Config, paymentCfg payment.PaymentCfg) Order {
	return &orderImpl{
		orderRepo: repository.CreateRepository[model.Order](cli, cfg),
		payment:   payment.NewPayment(paymentCfg),
	}
}

func (os *orderImpl) CreateOrder(pr model.PaymentRequest, paymentType payment.PaymentType) (any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	pr.Order.CreatedBy = pr.Order.Cart.UserID
	pr.Order.CreatedAt = time.Now()
	pr.Order.Status = model.OrderStatusPending

	paymentService, err := os.payment.GetService(paymentType)
	if paymentService == nil {
		return nil, err
	}

	psResp, err := paymentService.InitPayment(pr, "USD")
	if err != nil {
		return nil, err
	}
	var resp interface{}

	if paymentType == payment.PaymentTypePayPal {
		pr.Order.OrderID = psResp.(payment.PaypalCreateOrderResp).ID
		resp = psResp
	} else if paymentType == payment.PaymentTypeStripe {
		ses := psResp.(*stripe.CheckoutSession)
		pr.Order.OrderID = ses.ID
		resp = ses.URL
	}

	return resp, os.orderRepo.Add(ctx, pr.Order)
}

func (os *orderImpl) Checkout(oid string, userId model.UserID, paymentType payment.PaymentType) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err := os.payment.CaptureOrder(oid)
	if err != nil {
		return err
	}

	order := model.Order{
		OrderID: oid,
	}

	os.orderRepo.Get(ctx, order.Key(), &order)
	order.Status = model.OrderStatusApprooved

	return os.orderRepo.Update(ctx, order)
}
