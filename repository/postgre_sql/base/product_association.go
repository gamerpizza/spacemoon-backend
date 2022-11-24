package base

import (
	"context"
	"moonspace/model"
	"moonspace/repository/postgre_sql"
	"moonspace/repository/types"

	"gorm.io/gorm"
)

type ProductAssociatedRepositoryImpl[T types.ProductAssociation] struct {
	cli  *gorm.DB
	base BaseRepository[T]
}

func NewProductAssociatedRepository[T types.ProductAssociation](cli *gorm.DB) types.ProductAssociatedRepository[T] {
	return &ProductAssociatedRepositoryImpl[T]{
		cli:  cli,
		base: NewBasePostgresRepository[T](cli),
	}
}

func (cr *ProductAssociatedRepositoryImpl[T]) Add(c context.Context, data T) error {
	return cr.base.Add(c, data)
}

func (cr *ProductAssociatedRepositoryImpl[T]) Get(c context.Context, key map[string]interface{}, result *T) error {
	return cr.base.Get(c, key, result)
}

func (cr *ProductAssociatedRepositoryImpl[T]) GetLimit(c context.Context, start, end uint64, result *[]T) error {
	return cr.base.GetLimit(c, start, end, result)
}

func (cr *ProductAssociatedRepositoryImpl[T]) Delete(c context.Context, data T) error {
	return cr.base.Delete(c, data)
}

func (cr *ProductAssociatedRepositoryImpl[T]) Update(c context.Context, data T) error {
	return cr.base.Update(c, data)
}

func (cr *ProductAssociatedRepositoryImpl[T]) AddProduct(ctx context.Context, c *T, p *model.Product) error {
	query, value, err := postgre_sql.CreateQueryFromkey((*c).Key())
	if err != nil {
		return err
	}

	err = cr.cli.Model(c).Where(query, value).Association("Products").Append(p)
	return err
}

func (cr *ProductAssociatedRepositoryImpl[T]) GetProduct(ctx context.Context, c *T, pid string, result *model.Product) error {
	query, value, err := postgre_sql.CreateQueryFromkey((*c).Key())
	if err != nil {
		return err
	}

	return cr.cli.Model(c).Where(query, value).Association("Products").Find(&result, "product_id = ?", pid)
}

func (cr *ProductAssociatedRepositoryImpl[T]) GetProductsLimit(ctx context.Context, c *T, start, end uint64) error {
	// NOOP
	return nil
}

func (cr *ProductAssociatedRepositoryImpl[T]) UpdateProduct(ctx context.Context, c *T, p *model.Product) error {
	query, value, err := postgre_sql.CreateQueryFromkey((*c).Key())
	if err != nil {
		return err
	}

	return cr.cli.Model(c).Where(query, value).Association("Products").Replace(&model.Product{ID: p.ID}, p)
}

func (cr *ProductAssociatedRepositoryImpl[T]) DeleteProduct(ctx context.Context, c *T, p *model.Product) error {
	query, value, err := postgre_sql.CreateQueryFromkey((*c).Key())
	if err != nil {
		return err
	}

	return cr.cli.Model(c).Where(query, value).Association("Products").Delete(p)
}
