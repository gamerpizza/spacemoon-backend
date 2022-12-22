package network

import (
	"github.com/google/uuid"
	"spacemoon/login"
	"spacemoon/network/post"
	"time"
)

func New(p Persistence, user login.UserName) Network {
	return poster{persistence: p, user: user}
}

func NewPost(caption post.Caption, author login.UserName, urls post.ContentURIS) post.Post {
	var id = post.Id(uuid.NewString())
	return post.Post{Caption: caption, Author: author, URLS: urls, Id: id, Created: time.Now()}
}

type Network interface {
	Post(caption post.Caption, content ...post.ContentURI) (post.Post, error) //content
	GetPosts() (post.Posts, error)
}

type Persistence interface {
	AddPost(post post.Post) error
	GetAllPosts() (post.Posts, error)
}

type poster struct {
	user        login.UserName
	persistence Persistence
}

func (p poster) GetPosts() (posts post.Posts, err error) {
	posts, err = p.persistence.GetAllPosts()
	return
}

func (p poster) Post(c post.Caption, content ...post.ContentURI) (post.Post, error) {
	pst := p.makePost(c, content)
	err := p.persistence.AddPost(pst)
	if err != nil {
		return post.Post{}, err
	}
	return pst, nil
}

func (p poster) makePost(c post.Caption, content []post.ContentURI) post.Post {
	urls := make(post.ContentURIS)
	for _, url := range content {
		urls[url] = true
	}

	pst := NewPost(c, p.user, urls)
	return pst
}
