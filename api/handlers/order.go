package handlers

import (
	"moonspace/model"
	"moonspace/service"
	"moonspace/service/payment"
	"moonspace/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreatePayPalOrder(s service.Order) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		createOrder(ctx, s, payment.PaymentTypePayPal)
	}
}

func PayPalCheckout(s service.Order) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		checkout(ctx, s, payment.PaymentTypePayPal)
	}
}

func CreateStripeOrder(s service.Order) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		createOrder(ctx, s, payment.PaymentTypeStripe)
	}
}

func StripeCheckout(s service.Order) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		checkout(ctx, s, payment.PaymentTypeStripe)
	}
}

func CreateGooglePayOrder(s service.Order) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		createOrder(ctx, s, payment.PaymentTypeGoogle)
	}
}

func createOrder(ctx *gin.Context, s service.Order, pt payment.PaymentType) {
	pr := model.PaymentRequest{}

	err := utils.DecodeRequestBody(ctx.Request.Body, &pr)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	resp, err := s.CreateOrder(pr, pt)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(200, resp)
}

func checkout(ctx *gin.Context, s service.Order, pt payment.PaymentType) {
	id := ctx.Param("orderId")

	if pt == payment.PaymentTypePayPal {
		err := s.Checkout(id, model.UserID("1"), pt)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, nil)
	}
}
