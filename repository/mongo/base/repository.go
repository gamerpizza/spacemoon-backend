package base

import (
	"context"
	"errors"
	"fmt"
	"moonspace/model"
	"moonspace/repository/mongo/types"

	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
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

func (pr *BaseRepositoryImpl[T]) Update(c context.Context, data T) error {
	col, err := pr.SwitchTable(&data)
	if err != nil {
		return err
	}

	fmt.Println(data.Key())
	_, err = col.Upsert(c, bson.M(data.Key()), data)
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
