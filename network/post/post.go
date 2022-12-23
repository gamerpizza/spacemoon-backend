// Package post contains the basic domain definition fot the Posts used in the network package
package post

import (
	"spacemoon/login"
	"time"
)

type Post struct {
	Caption Caption        `json:"caption"`
	Author  login.UserName `json:"author"`
	URLS    ContentURIS    `json:"urls"`
	Id      Id             `json:"id"`
	Created time.Time      `json:"created"`
	Likes   Likes          `json:"likes"`
}

func (p *Post) GetCaption() Caption {
	return p.Caption
}

func (p *Post) GetId() Id {
	return p.Id
}

func (p *Post) Comments() Comments {
	return Comments{}
}

func (p *Post) Content() Content {
	return Content{URLS: p.URLS}
}

func (p *Post) GetAuthor() login.UserName {
	return p.Author
}

func (p *Post) AddLike(u login.UserName) {
	if p.Likes == nil {
		p.Likes = make(Likes)
	}
	p.Likes[u] = true
}

func (p *Post) RemoveLike(u login.UserName) {
	delete(p.Likes, u)
}

func (p *Post) GetLikes() Likes {
	return p.Likes
}

// Caption is the text added to a post
type Caption string

// Id is a unique UUID string to identify every Post
type Id string
type Posts map[Id]Post
type Likes map[login.UserName]bool
