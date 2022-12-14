package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"spacemoon/login"
	"spacemoon/product"
	"spacemoon/product/category"
	"spacemoon/product/ratings"
)

type Persistence interface {
	login.Persistence
	product.Persistence
	ratings.Persistence
	category.Persistence
	Close() error
}

func GetPersistence(ctx context.Context) (Persistence, error) {
	const projectID = "global-pagoda-368419"
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("could not create firestore client: %w", err)
	}

	return &fireStorePersistence{ctx: ctx, storage: client}, nil
}

type fireStorePersistence struct {
	storage *firestore.Client
	ctx     context.Context
}

func (p *fireStorePersistence) Close() error {
	err := p.storage.Close()
	if err != nil {
		return err
	}
	return nil
}

func (p *fireStorePersistence) GetCategories() category.Categories {
	return nil
}

func (p *fireStorePersistence) SaveCategory(dto category.DTO) {
	//TODO implement me
	panic("implement me")
}

func (p *fireStorePersistence) DeleteCategory(name category.Name) {
	//TODO implement me
	panic("implement me")
}

func (p *fireStorePersistence) ReadRating(_ product.Id) ratings.Rating {
	return ratings.Rating{}
}

func (p *fireStorePersistence) SaveRating(id product.Id, rating ratings.Rating) {
	//TODO implement me
	panic("implement me")
}

const loginCollection = "login"
const loginTokensCollection = "login-tokens"
const productCollection = "products"
