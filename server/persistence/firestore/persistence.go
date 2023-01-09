package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"fmt"
	"spacemoon/login"
	"spacemoon/network"
	"spacemoon/network/message"
	"spacemoon/network/post"
	"spacemoon/network/post/comment"
	"spacemoon/network/profile"
	"spacemoon/product"
	"spacemoon/product/category"
	"spacemoon/product/ratings"
	"strings"
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

func (p *fireStorePersistence) GetProfile(id profile.Id) (profile.Profile, error) {
	collection := p.storage.Collection(profilesCollection)
	doc, err := collection.Doc(string(id)).Get(p.ctx)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			return profile.Profile{}, NotFoundError
		}
		return profile.Profile{}, fmt.Errorf("could not read from firestore: %w", err)
	}
	pr := profile.Profile{}
	err = doc.DataTo(&pr)
	if err != nil {
		return profile.Profile{}, fmt.Errorf("could not parse data from persistence into profile: %w", err)
	}
	return pr, nil
}

func (p *fireStorePersistence) SaveProfile(pr profile.Profile) error {
	collection := p.storage.Collection(profilesCollection)
	_, err := collection.Doc(string(pr.Id)).Set(p.ctx, pr)
	if err != nil {
		return fmt.Errorf("could not save to firestore: %w", err)
	}
	return nil
}

func (p *fireStorePersistence) DeletePost(id post.Id) error {
	collection := p.storage.Collection(postsCollection)
	_, err := collection.Doc(string(id)).Delete(p.ctx)
	if err != nil {
		return fmt.Errorf("could not delete from collection: %w", err)
	}
	return nil
}

func (p *fireStorePersistence) AddPost(post post.Post) error {
	collection := p.storage.Collection(postsCollection)
	_, err := collection.Doc(string(post.GetId())).Set(p.ctx, post)
	if err != nil {
		return fmt.Errorf("could not write to collection: %w", err)
	}
	return nil
}

func (p *fireStorePersistence) GetAllPosts() (post.Posts, error) {
	collection := p.storage.Collection(postsCollection)
	documents, err := collection.Documents(p.ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("could not write to collection: %w", err)
	}

	posts := post.Posts{}

	for _, document := range documents {
		var pst post.Post
		err = document.DataTo(&pst)
		if err != nil {
			return nil, fmt.Errorf("could parse document: %w", err)
		}
		posts[pst.GetId()] = pst
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

const projectID = "global-pagoda-368419"
const productCollection = "products"
const postsCollection = "posts"
const profilesCollection = "profiles"
const messagesCollection = "direct-messages"
const commentsCollection = "post-comments"
const commentsSubCollection = "comments"

var NotFoundError = errors.New("not found")
