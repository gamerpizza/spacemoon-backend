package base

import (
	"context"
	"moonspace/model"
	"moonspace/repository/types"

	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type ProductAssociatedRepositoryImpl[T types.ProductAssociation] struct {
	cli  *qmgo.QmgoClient
	base BaseRepository[T]
}

func NewProductAssociatedRepository[T types.ProductAssociation](cli *qmgo.QmgoClient, dbName string) types.ProductAssociatedRepository[T] {
	return &ProductAssociatedRepositoryImpl[T]{
		cli:  cli,
		base: NewBaseMongoRepository[T](cli, dbName),
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

func (cr *ProductAssociatedRepositoryImpl[T]) SwitchTable(data *T) (*qmgo.Collection, error) {
	return cr.base.SwitchTable(data)
}

func (cr *ProductAssociatedRepositoryImpl[T]) AddProduct(ctx context.Context, c *T, p *model.Product) error {
	coll, err := cr.SwitchTable(c)
	if err != nil {
		return err
	}

	updateQuery := bson.M{"$push": bson.M{"products": p}}

	return coll.UpdateOne(ctx, (*c).Key(), updateQuery)
}

func (cr *ProductAssociatedRepositoryImpl[T]) GetProduct(ctx context.Context, c *T, pid string, p *model.Product) error {
	coll, err := cr.SwitchTable(c)
	if err != nil {
		return err
	}

	entityQuery := (bson.M)((*c).Key())
	entityQuery["products._id"] = pid

	return coll.Find(ctx, entityQuery).One(p)
}

func (cr *ProductAssociatedRepositoryImpl[T]) GetProductsLimit(ctx context.Context, c *T, start, end uint64) error {
	coll, err := cr.SwitchTable(c)
	if err != nil {
		return err
	}

	entityQuery := (bson.M)((*c).Key())
	return coll.Find(ctx, entityQuery).
		Select(bson.M{"products": bson.M{"$slice": []uint64{start, end}}}).One(c)
}

func (cr *ProductAssociatedRepositoryImpl[T]) UpdateProduct(ctx context.Context, c *T, p *model.Product) error {
	coll, err := cr.SwitchTable(c)
	if err != nil {
		return err
	}

	updateQuery := bson.M{"products.product_id": p.ID, "$set": *p}
	return coll.UpdateOne(ctx, (*c).Key(), updateQuery)
}

func (cr *ProductAssociatedRepositoryImpl[T]) DeleteProduct(ctx context.Context, c *T, p *model.Product) error {
	coll, err := cr.SwitchTable(c)
	if err != nil {
		return err
	}

	updateQuery := bson.M{"$pull": bson.M{"products": p.Key()}}
	return coll.UpdateOne(ctx, (*c).Key(), updateQuery)
}
