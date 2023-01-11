package firestore

import (
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"spacemoon/network/post"
	"strings"
)

func (p *fireStorePersistence) DeletePost(id post.Id) error {
	collection := p.storage.Collection(postsCollection)
	_, err := collection.Doc(string(id)).Delete(p.ctx)
	if err != nil {
		return fmt.Errorf("could not delete from collection: %w", err)
	}
	return nil
}

func (p *fireStorePersistence) AddPost(post post.Post) error {
	if strings.TrimSpace(string(post.GetId())) == "" {
		return EmptyIdOnPostError
	}
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

func (p *fireStorePersistence) CheckIfPostExists(id post.Id) (bool, error) {
	collection := p.storage.Collection(postsCollection)
	_, err := collection.Doc(string(id)).Get(p.ctx)
	if err != nil && status.Code(err) == codes.NotFound {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("error checking post existence: %w", err)
	}
	return true, nil
}

func (p *fireStorePersistence) GetPost(id post.Id) (post.Post, error) {
	collection := p.storage.Collection(postsCollection)
	doc, err := collection.Doc(string(id)).Get(p.ctx)
	if err != nil {
		return post.Post{}, fmt.Errorf("could not get saved post: %w", err)
	}
	var pst post.Post
	err = doc.DataTo(&pst)
	if err != nil {
		return post.Post{}, fmt.Errorf("could not parse saved post: %w", err)
	}
	return pst, nil
}

var EmptyIdOnPostError = errors.New("post cannot have an empty id")
