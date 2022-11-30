package base

import (
	"context"
	"errors"
	"moonspace/model"
	"moonspace/repository/postgre_sql"
	"moonspace/repository/postgre_sql/types"

	"gorm.io/gorm"
)

var (
	errUnsupportedRepository = errors.New("unsupported repository type")
	errConvertingData        = errors.New("error converting query result")
)

type BaseRepository[T model.Entity] interface {
	types.PostgreRepository[T]
}

type BaseRepositoryImpl[T model.Entity] struct {
	cli *gorm.DB
}

func (pr *BaseRepositoryImpl[T]) Add(c context.Context, data T) error {
	tx := pr.cli.Create(&data)
	if tx.Error != nil {
		return tx.Error
	}

	return tx.Save(&data).Error
}

func (pr *BaseRepositoryImpl[T]) Get(c context.Context, key map[string]interface{}, result *T) error {
	query, val, err := postgre_sql.CreateQueryFromkey(key)
	if err != nil {
		return err
	}

	tx := pr.cli.Where(query, val).First(&result)
	return tx.Error
}

func (pr *BaseRepositoryImpl[T]) GetLimit(c context.Context, start, end uint64, result *[]T) error {
	// NOOP For now
	return nil
}

func (pr *BaseRepositoryImpl[T]) Delete(c context.Context, data T) error {
	query, val, err := createQuery(data)
	if err != nil {
		return err
	}

	return pr.cli.Where(query, val).Delete(&data).Error
}

func (pr *BaseRepositoryImpl[T]) Update(c context.Context, id string, data T) error {
	query, val, err := createQuery(data)
	if err != nil {
		return err
	}

	tx := pr.cli.Model(data).Where(query, val...).Updates(&data)

	return tx.Error
}

func (pr *BaseRepositoryImpl[T]) GetProductLimit(c context.Context, cid string, start, end uint64, result *[]T) error {
	// NOOP For now
	return nil
}

func (pr *BaseRepositoryImpl[T]) DeleteProduct(c context.Context, cid string, data T) error {
	// NOOP For now
	return nil
}

func (pr *BaseRepositoryImpl[T]) UpdateProduct(c context.Context, cid string, id string, data T) error {
	// NOOP For now
	return nil
}

func createQuery[T model.Entity](data T) (string, []interface{}, error) {
	key := data.Key()

	return postgre_sql.CreateQueryFromkey(key)
}

func NewBasePostgresRepository[T model.Entity](
	cli *gorm.DB,
) types.PostgreRepository[T] {
	return &BaseRepositoryImpl[T]{
		cli,
	}
}
