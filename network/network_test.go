package network

import (
	"spacemoon/network/post"
	"testing"
)

type mockPersistence struct {
	posts map[post.Id]post.Post
}

func (m *mockPersistence) AddPost(p post.Post) error {
	if m.posts == nil {
		m.posts = make(map[post.Id]post.Post)
	}
	m.posts[p.GetId()] = p
	return nil
}

func (m *mockPersistence) GetAllPosts() (post.Posts, error) {
	return m.posts, nil
}

func (m *mockPersistence) DeletePost(id post.Id) error {
	//TODO implement me
	panic("implement me")
}

func TestNetwork(t *testing.T) {
	const user = "Edgar Allan post"
	var n PostManager = NewPostManager(&mockPersistence{}, user)
	const caption = "something"
	const imageUrl1 = "image-url-1"
	const imageUrl2 = "image-url-2"
	const videoUrl1 = "video-url-1"
	const videoUrl2 = "video-url-2"
	var p, _ = n.Post(caption, imageUrl1, imageUrl2, videoUrl1, videoUrl2)
	if p.GetCaption() != caption {
		t.Fatalf("invalid Caption (%s), expected: %s", p.GetCaption(), caption)
	}
	if p.GetAuthor() != user {
		t.Fatalf("invalid GetAuthor (%s), expected: %s", p.GetAuthor(), user)
	}
	var c post.Content = p.Content()
	var urls post.ContentURIS = c.GetURLS()
	if urls.Is(imageUrl1).NotPresent() {
		t.Fatalf("(%s) url not found", imageUrl1)
	}
	if urls.Is(imageUrl2).NotPresent() {
		t.Fatalf("(%s) url not found", imageUrl2)
	}
	if urls.Is(videoUrl1).NotPresent() {
		t.Fatalf("(%s) url not found", videoUrl1)
	}
	if urls.Is(videoUrl2).NotPresent() {
		t.Fatalf("(%s) url not found", videoUrl2)
	}

	var _ post.Id = p.GetId()
	var _ post.Comments = p.Comments()
}

func TestNetwork_GetPosts(t *testing.T) {
	const user = "Edgar Allan post"
	var n PostManager = NewPostManager(&mockPersistence{}, user)
	const caption1 = "something"
	const caption2 = "something other"
	post1, _ := n.Post(caption1)
	post2, _ := n.Post(caption2)
	var retrievedPosts, _ = n.GetPosts()
	if retrieved, exists := retrievedPosts[post1.GetId()]; !exists || retrieved.GetCaption() != caption1 {
		t.Fatal("posted post not found")
	}
	if retrieved, exists := retrievedPosts[post2.GetId()]; !exists || retrieved.GetCaption() != caption2 {
		t.Fatal("posted post not found")
	}

}
