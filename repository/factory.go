package repository

import (
	"moonspace/model"
	mongo_base "moonspace/repository/mongo/base"
	postgres_base "moonspace/repository/postgre_sql/base"
	"moonspace/repository/types"

	"github.com/qiniu/qmgo"
	"gorm.io/gorm"
)

func CreateRepository[T model.Entity](cli any, cfg types.Config) types.Repository[T] {
	switch cfg.Type {
	case types.Mongo:
		return mongo_base.NewBaseMongoRepository[T](cli.(*qmgo.QmgoClient), cfg.Database)
	case types.Postgres:
		return postgres_base.NewBasePostgresRepository[T](cli.(*gorm.DB))
	default:
		panic("unsupported database type: " + cfg.Type)
	}
}

func CreateProductRepository[T types.ProductAssociation](cli any, cfg types.Config) types.ProductAssociatedRepository[T] {
	switch cfg.Type {
	case types.Mongo:
		return mongo_base.NewProductAssociatedRepository[T](cli.(*qmgo.QmgoClient), cfg.Database)
	case types.Postgres:
		return postgres_base.NewProductAssociatedRepository[T](cli.(*gorm.DB))
	default:
		panic("unsupported database type: " + cfg.Type)
	}
}
