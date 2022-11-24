package types

import (
	"context"
	"moonspace/model"
	"time"
)

type Transaction func(sessCtx context.Context) (interface{}, error)

type Repository[T model.Entity] interface {
	Add(c context.Context, data T) error
	Get(c context.Context, key map[string]interface{}, result *T) error
	GetLimit(c context.Context, start, end uint64, result *[]T) error
	Delete(c context.Context, data T) error
	Update(c context.Context, data T) error
}

type ProductAssociation interface {
	model.Cart | model.Category
	model.Entity
	model.ProductEntity
}

type ProductAssociatedRepository[T ProductAssociation] interface {
	AddProduct(ctx context.Context, c *T, p *model.Product) error
	DeleteProduct(ctx context.Context, c *T, p *model.Product) error
	GetProduct(ctx context.Context, c *T, pid string, p *model.Product) error
	GetProductsLimit(ctx context.Context, c *T, start, end uint64) error
	UpdateProduct(ctx context.Context, c *T, p *model.Product) error
	Repository[T]
}

const (
	Timeout = time.Second * 10
)
