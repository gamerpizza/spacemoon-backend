package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"spacemoon/login"
	"spacemoon/network"
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

func (p *fireStorePersistence) AddPost(post network.Post) error {
	collection := p.storage.Collection(postsCollection)
	_, err := collection.Doc(string(post.GetId())).Set(p.ctx, post)
	if err != nil {
		return fmt.Errorf("could not write to collection: %w", err)
	}
	return nil
}

func (p *fireStorePersistence) GetAllPosts() (network.Posts, error) {
	collection := p.storage.Collection(postsCollection)
	documents, err := collection.Documents(p.ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("could not write to collection: %w", err)
	}

	posts := network.Posts{}

	for _, document := range documents {
		var post network.Post
		err = document.DataTo(&post)
		if err != nil {
			return nil, fmt.Errorf("could parse document: %w", err)
		}
		posts[post.GetId()] = post
	}
	return posts, nil
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
const productCollection = "products"
const postsCollection = "posts"
