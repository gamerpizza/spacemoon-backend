package post

import (
	"reflect"
	"spacemoon/login"
	"testing"
	"time"
)

func TestPoster(t *testing.T) {
	const caption = "test post"
	const author = "test author"
	const id = "test id"
	const uri = "test uri"
	var p Post = Post{
		Caption: caption,
		Author:  author,
		URLS:    ContentURIS{uri: true},
		Id:      id,
		Created: time.Time{},
	}
	if p.GetCaption() != caption {
		t.Fatal("bad caption")
	}
	if p.GetAuthor() != author {
		t.Fatal("bad author")
	}
	if p.GetId() != id {
		t.Fatal("bad id")
	}
	if _, exists := p.Content().GetURLS()[uri]; !exists {
		t.Fatal("did not find the expected uri")
	}
	if p.Content().GetURLS().Is(uri).NotPresent() {
		t.Fatal("Is Not Present method not working")
	}
	var u login.UserName

	p.AddLike(u)
	if _, exists := p.Likes[u]; !exists {
		t.Fatal("like not added")
	}
	p.RemoveLike(u)
	if _, exists := p.Likes[u]; exists {
		t.Fatal("like not removed")
	}
	if !reflect.DeepEqual(p.Likes, p.GetLikes()) {
		t.Fatal("GetLikes did not return Likes")
	}
}
