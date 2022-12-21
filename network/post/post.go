package post

import (
	"spacemoon/login"
	"time"
)

type Post struct {
	Caption Caption        `json:"caption"`
	Author  login.UserName `json:"author"`
	URLS    ContentURLS    `json:"urls"`
	Id      Id             `json:"id"`
	Created time.Time      `json:"created"`
}

func (p Post) GetCaption() Caption {
	return p.Caption
}

func (p Post) GetId() Id {
	return p.Id
}

func (p Post) Comments() Comments {
	return Comments{}
}

func (p Post) Content() Content {
	return Content{URLS: p.URLS}
}

func (p Post) GetAuthor() login.UserName {
	return p.Author
}

type Caption string
type Id string
type Posts map[Id]Post
