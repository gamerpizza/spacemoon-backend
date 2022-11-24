package repository_test

import (
	"context"
	"moonspace/api"
	"moonspace/model"
	"moonspace/repository"
	mongo "moonspace/repository/mongo/base"
	"moonspace/repository/types"
	"moonspace/test_helpers"
	"testing"
	"time"

	"github.com/qiniu/qmgo"
	"github.com/stretchr/testify/assert"
)

func Test_Mongo(t *testing.T) {
	t.Parallel()

	cfg := api.Config{
		DB: &types.Config{
			Url:      "mongodb://localhost:27017",
			Database: "spacemoon_test",
			Type:     types.Mongo,
		},
	}

	cli := repository.MakeRepositoryClient(*cfg.DB)

	productRepo := repository.CreateRepository[model.Product](cli, *cfg.DB)
	cartRepo := mongo.NewProductAssociatedRepository[model.Cart](cli.(*qmgo.QmgoClient), cfg.DB.Database)
	categoryRepo := mongo.NewProductAssociatedRepository[model.Category](cli.(*qmgo.QmgoClient), cfg.DB.Database)
	orderRepo := repository.CreateRepository[model.Order](cli, *cfg.DB)
	testCRUD(t, productRepo, cartRepo, categoryRepo, orderRepo, true)
}

func Test_Postgres(t *testing.T) {
	// t.Parallel()

	// cfg := api.Config{
	// 	DB: &types.Config{
	// 		Url:      "host=localhost user=postgres password=admin dbname=spacemoon-test port=5432 sslmode=disable",
	// 		Database: "spacemoon",
	// 		Type:     types.Postgres,
	// 	},
	// }

	// cli := repository.MakeRepositoryClient(*cfg.DB)

	// productRepo := repository.CreateRepository[model.Product](cli, *cfg.DB)
	// cartRepo := base.NewProductAssociatedRepository[model.Cart](cli.(*gorm.DB))
	// categoryRepo := base.NewProductAssociatedRepository[model.Category](cli.(*gorm.DB))
	// testCRUD(t, productRepo, cartRepo, categoryRepo)
}

func testCRUD(
	t *testing.T,
	productRepo types.Repository[model.Product],
	cartRepo types.ProductAssociatedRepository[model.Cart],
	categoryRepo types.ProductAssociatedRepository[model.Category],
	orderRepo types.Repository[model.Order],
	isMongo bool,
) {
	prod, cart, category, order, _ := test_helpers.CreateTestData()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// START CATEGORY
	assert.NoError(t, categoryRepo.Add(ctx, category))
	myCat := model.Category{}
	assert.NoError(t, categoryRepo.Get(ctx, category.Key(), &myCat))
	assert.NotNil(t, myCat)

	for _, v := range prod {
		assert.NoError(t, categoryRepo.AddProduct(ctx, &category, &v))
	}
	// END CATEGORY

	// START CART
	assert.NoError(t, cartRepo.Add(ctx, cart))
	for _, v := range prod {
		assert.NoError(t, cartRepo.AddProduct(ctx, &cart, &v))
	}
	myCart := model.Cart{}
	assert.NoError(t, cartRepo.Get(ctx, cart.Key(), &myCart))
	assert.NotNil(t, myCart)
	assert.Equal(t, 3, len(myCart.Products))
	// END CART

	// START ORDER
	order.Cart = myCart
	assert.NoError(t, orderRepo.Add(ctx, order))
	o := model.Order{}
	assert.NoError(t, orderRepo.Get(ctx, order.Key(), &o))
	assert.NotNil(t, o)
	assert.Equal(t, 3, len(o.Cart.Products))
	// END ORDER

	if !isMongo {
		p1 := model.Product{}
		err := productRepo.Get(ctx, prod[0].Key(), &p1)
		assert.NoError(t, err)
		assert.NotNil(t, p1)

		p2 := model.Product{}
		err = productRepo.Get(ctx, prod[1].Key(), &p2)
		assert.NoError(t, err)
		assert.NotNil(t, p1)

		p3 := model.Product{}
		err = productRepo.Get(ctx, prod[2].Key(), &p3)
		assert.NoError(t, err)
		assert.NotNil(t, p1)

		assert.NoError(t, productRepo.Delete(ctx, p1))
		assert.NoError(t, productRepo.Delete(ctx, p2))
		assert.NoError(t, productRepo.Delete(ctx, p3))
	}

	// CLEANUP
	assert.NoError(t, categoryRepo.Delete(ctx, category))
	assert.NoError(t, cartRepo.Delete(ctx, cart))
	assert.NoError(t, orderRepo.Delete(ctx, o))
}
