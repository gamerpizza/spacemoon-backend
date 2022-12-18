package network

import (
	"github.com/google/uuid"
	"spacemoon/login"
)

type Network interface {
	Post(caption PostCaption, content ...PostContentURL) (Post, error)
	GetPosts() (Posts, error)
}

type Persistence interface {
	AddPost(post Post) error
	GetAllPosts() (Posts, error)
}

func New(p Persistence, user login.UserName) Network {
	return poster{persistence: p, user: user}
}

type poster struct {
	user        login.UserName
	persistence Persistence
}

func (p poster) GetPosts() (posts Posts, err error) {
	posts, err = p.persistence.GetAllPosts()
	return
}

func (p poster) Post(c PostCaption, content ...PostContentURL) (Post, error) {
	post := p.makePost(c, content)
	err := p.persistence.AddPost(post)
	if err != nil {
		return Post{}, err
	}
	return post, nil
}

func (p poster) makePost(c PostCaption, content []PostContentURL) Post {
	urls := make(PostContentURLS)
	for _, url := range content {
		urls[url] = nil
	}

	post := NewPost(c, p.user, urls)
	return post
}

func NewPost(caption PostCaption, author login.UserName, urls PostContentURLS) Post {
	var id = PostId(uuid.NewString())
	return Post{Caption: caption, Author: author, URLS: urls, Id: id}
}

type PostContentURL string
