package payment

import (
	"fmt"
	"moonspace/model"
	"strings"

	"github.com/stripe/stripe-go/v73"
	"github.com/stripe/stripe-go/v73/checkout/session"
)

type Stripe struct {
	cfg StripeCfg
}

func NewStripeClient(cfg StripeCfg) Payment {
	stripe.Key = cfg.APIKey

	fmt.Println("STRIPE: ", cfg)
	return &Stripe{
		cfg,
	}
}

// InitPayment creates a payment
func (s *Stripe) InitPayment(pr model.PaymentRequest, currency string) (any, error) {
	curr := stripe.String(strings.ToLower(currency))

	ses := &stripe.CheckoutSessionParams{
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		CancelURL:  stripe.String(s.cfg.CancelUrl),
		SuccessURL: stripe.String(s.cfg.SuccessUrl),
		Locale:     stripe.String("en"),
		Currency:   curr,
		LineItems:  nil,
	}
	items := make([]*stripe.CheckoutSessionLineItemParams, 0)

	for _, p := range pr.Order.Cart.Products {
		items = append(items, &stripe.CheckoutSessionLineItemParams{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency: curr,
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Description: stripe.String(p.Description),
					Name:        stripe.String(p.Name),
					Images:      stripe.StringSlice([]string{p.Image}),
				},
				UnitAmountDecimal: stripe.Float64(p.Price * 100),
			},
			Quantity: stripe.Int64(1),
		})
	}

	ses.LineItems = items

	newSes, err := session.New(ses)
	if err != nil {
		return newSes.CancelURL, err
	}

	return newSes, nil
}

// CaptureOrder will act as a webhook for Stripe.
// It verifies if the payment is successfull or not.
// If not, it will delete the order, and restock the product.
func (s *Stripe) CaptureOrder(id string) (any, error) {
	return nil, nil
}
