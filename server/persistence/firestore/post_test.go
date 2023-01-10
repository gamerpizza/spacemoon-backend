package firestore

import (
	"context"
	"spacemoon/network"
	"spacemoon/network/post"
	"testing"
)

func TestFireStorePersistence_CheckIfPostExists(t *testing.T) {
	var p network.Persistence
	var err error
	p, err = GetPersistence(context.TODO())
	var newPost post.Post = post.New("test-post", "testMan", nil)
	err = p.AddPost(newPost)
	if err != nil {
		t.Fatal("could not save post")
	}
	exists, err := p.CheckIfPostExists(newPost.Id)
	if err != nil {
		t.Errorf("could not check post: %s", err.Error())
	}
	if !exists {
		t.Error("post not found")
	}
	err = p.DeletePost(newPost.Id)
	if err != nil {
		t.Error("could not delete post")
	}
	exists, err = p.CheckIfPostExists(newPost.Id)
	if err != nil {
		t.Error("could not check post")
	}
	if exists {
		t.Error("post not deleted")
	}
}
