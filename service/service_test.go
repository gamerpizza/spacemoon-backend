package service_test

import (
	"context"
	"moonspace/api"
	"moonspace/model"
	"moonspace/repository"
	"moonspace/repository/types"
	"moonspace/service"
	"moonspace/service/payment"
	"moonspace/test_helpers"
	"testing"
	"time"

	"github.com/qiniu/qmgo"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func Test_Service_Mongo(t *testing.T) {
	cfg := api.Config{
		DB: &types.Config{
			Url:      "mongodb://localhost:27017",
			Database: "spacemoon_test",
			Type:     types.Mongo,
		},
		Payment: &payment.PaymentCfg{
			PayPal: payment.PaypalCfg{
				ClientID:     "test",
				ClientSecret: "test",
				URL:          "https://api-m.sandbox.paypal.com",
			},
			Stripe: payment.StripeCfg{
				PublicKey: "asdf",
				APIKey:    "zxcv",
			},
		},
	}
	cli := repository.MakeRepositoryClient(*cfg.DB)

	s := service.NewService(cli, *cfg.DB, *cfg.Payment)
	prod, _, category, order, _ := test_helpers.CreateTestData()
	productNumber := len(prod)

	assert.NoError(t, s.Category.Create(category))

	for _, p := range prod {
		assert.NoError(t, s.Product.Create(p))
	}

	catDB, err := s.Category.Get(category.ID)
	assert.NoError(t, err)
	assert.NotNil(t, catDB)
	assert.Equal(t, productNumber, len(catDB.Products))

	for _, pr := range prod {
		p, err := s.Product.Get(category.ID, pr.ID)
		assert.NoError(t, err)
		assert.NotNil(t, p)
	}

	pr := model.PaymentRequest{
		Order: order,
	}
	assert.NoError(t, err)
	resp, err := s.Order.CreateOrder(pr, payment.PaymentTypePayPal)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	for _, pr := range prod {
		assert.NoError(t, s.Product.Delete(pr.CategoryID, pr.ID))
	}

	assert.NoError(t, s.Category.Delete(category.ID))
	qmgoCli := cli.(*qmgo.QmgoClient)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	qmgoCli.Client.Database(cfg.DB.Database).Collection("order").RemoveAll(ctx, bson.M{})
}
