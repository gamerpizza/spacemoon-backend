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
	Update(c context.Context, id string, data T) error
	GetProductLimit(c context.Context, cid string, start, end uint64, result *[]T) error
	DeleteProduct(c context.Context, cid string, data T) error
	UpdateProduct(c context.Context, cid string, id string, data T) error
}

type ProductAssociation interface {
	model.Cart | model.Category
	model.Entity
	model.ProductEntity
}

const (
	Timeout = time.Second * 10
)
