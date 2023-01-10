package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"fmt"
	"spacemoon/login"
	"spacemoon/network"
	"spacemoon/network/message"
	"spacemoon/network/post/comment"
	"spacemoon/network/profile"
	"spacemoon/product"
	"spacemoon/product/category"
	"spacemoon/product/ratings"
)

type Persistence interface {
	login.Persistence
	product.Persistence
	ratings.Persistence
	category.Persistence
	network.Persistence
	profile.Persistence
	message.Persistence
	comment.Persistence
	Close() error
}

func GetPersistence(ctx context.Context) (Persistence, error) {
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

const projectID = "global-pagoda-368419"
const productCollection = "products"
const postsCollection = "posts"
const profilesCollection = "profiles"
const messagesCollection = "direct-messages"
const commentsCollection = "post-comments"
const commentsSubCollection = "comments"

var NotFoundError = errors.New("not found")
