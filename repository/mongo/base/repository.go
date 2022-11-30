package base

import (
	"context"
	"errors"
	"moonspace/model"
	"moonspace/repository/mongo/types"

	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var errUnsupportedRepository = errors.New("cannot switch table; unknown data type")

type BaseRepository[T model.Entity] interface {
	types.MongoRepository[T]
}

type BaseRepositoryImpl[T model.Entity] struct {
	cli            *qmgo.Database
	transactionCli *qmgo.QmgoClient
}

func (pr *BaseRepositoryImpl[T]) Add(c context.Context, data T) error {
	col, err := pr.SwitchTable(&data)
	if err != nil {
		return err
	}

	_, err = col.InsertOne(c, data)

	if err != nil {
		return err
	}

	return nil
}

func (pr *BaseRepositoryImpl[T]) Get(c context.Context, key map[string]interface{}, result *T) error {
	col, err := pr.SwitchTable(result)
	if err != nil {
		return err
	}

	return col.Find(c, key).One(result)
}

func (pr *BaseRepositoryImpl[T]) GetLimit(c context.Context, start, end uint64, result *[]T) error {
	var t T
	col, err := pr.SwitchTable(&t)
	if err != nil {
		return err
	}

	return col.Find(c, bson.M{}).Skip(int64(start)).Limit(int64(end - start)).All(result)
}

func (pr *BaseRepositoryImpl[T]) Delete(c context.Context, data T) error {
	col, err := pr.SwitchTable(&data)
	if err != nil {
		return err
	}

	return col.Remove(c, bson.M(data.Key()))
}

func (pr *BaseRepositoryImpl[T]) Update(c context.Context, id string, data T) error {
	col, err := pr.SwitchTable(&data)
	if err != nil {
		return err
	}

	oid, _ := primitive.ObjectIDFromHex(id)
	_, err = col.Upsert(c, bson.M{"_id": oid}, data)
	if err != nil {
		return err
	}

	return nil
}

func (pr *BaseRepositoryImpl[T]) GetProductLimit(c context.Context, cid string, start, end uint64, result *[]T) error {
	var t T
	col, err := pr.SwitchTable(&t)
	if err != nil {
		return err
	}

	entityQuery := bson.M{"category_id": cid}

	return col.Find(c, entityQuery).Skip(int64(start)).Limit(int64(end - start)).All(result)
}

func (pr *BaseRepositoryImpl[T]) DeleteProduct(c context.Context, cid string, data T) error {
	col, err := pr.SwitchTable(&data)
	if err != nil {
		return err
	}
	entityQuery := (bson.M)(data.Key())
	entityQuery["category_id"] = cid
	return col.Remove(c, entityQuery)
}

func (pr *BaseRepositoryImpl[T]) UpdateProduct(c context.Context, cid string, id string, data T) error {
	col, err := pr.SwitchTable(&data)
	if err != nil {
		return err
	}

	oid, _ := primitive.ObjectIDFromHex(id)
	entityQuery := bson.M{"_id": oid}
	entityQuery["category_id"] = cid

	_, err = col.Upsert(c, entityQuery, data)
	if err != nil {
		return err
	}

	return nil
}

func (pr *BaseRepositoryImpl[T]) SwitchTable(data *T) (*qmgo.Collection, error) {
	switch any(*data).(type) {
	case model.Product:
		return pr.cli.Collection("product"), nil
	case model.Order:
		return pr.cli.Collection("order"), nil
	case model.Cart:
		return pr.cli.Collection("cart"), nil
	case model.Category:
		return pr.cli.Collection("category"), nil
	default:
		return nil, errUnsupportedRepository
	}
}

func NewBaseMongoRepository[T model.Entity](
	cli *qmgo.QmgoClient,
	dbName string,
) types.MongoRepository[T] {
	return &BaseRepositoryImpl[T]{
		cli:            cli.Client.Database(dbName),
		transactionCli: cli,
	}
}