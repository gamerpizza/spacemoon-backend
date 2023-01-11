// Package network is the basis of spacemoon's social network. It is tasked with managing the creating and fetching
// of posts and media information associated to them
package network

import (
	"spacemoon/login"
	"spacemoon/network/post"
)

// NewPostManager is the only way to externally instantiate new PostManager from this package,
// so that any user of this package has access to the PostManager only through its interface
func NewPostManager(p Persistence, user login.UserName) PostManager {
	return postManager{persistence: p, user: user}
}

// PostManager gets and creates post.Posts
type PostManager interface {
	// Post creates a new post
	Post(caption post.Caption, content ...string) (post.Post, error)
	// GetPosts retrieves existing posts
	GetPosts() (post.Posts, error)
}

// Persistence defines how the network package persists information
type Persistence interface {
	AddPost(post post.Post) error
	GetAllPosts() (post.Posts, error)
	DeletePost(post.Id) error
	CheckIfPostExists(post.Id) (bool, error)
}

type postManager struct {
	user        login.UserName
	persistence Persistence
}

func (p postManager) GetPosts() (posts post.Posts, err error) {
	posts, err = p.persistence.GetAllPosts()
	return
}

func (p postManager) Post(c post.Caption, content ...string) (post.Post, error) {
	pst := p.makePost(c, content)
	err := p.persistence.AddPost(pst)
	if err != nil {
		return post.Post{}, err
	}
	return pst, nil
}

func (p postManager) makePost(c post.Caption, content []string) post.Post {
	urls := make(post.ContentURIS)
	for _, url := range content {
		urls[url] = true
	}

	pst := post.New(c, p.user, urls)
	return pst
}
