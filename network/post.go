package network

import "spacemoon/login"

type Post struct {
	Caption PostCaption     `json:"caption"`
	Author  login.UserName  `json:"author"`
	URLS    PostContentURLS `json:"urls"`
	Id      PostId          `json:"id"`
}

func (p Post) GetCaption() PostCaption {
	return p.Caption
}

func (p Post) GetId() PostId {
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

type PostCaption string
type PostId string
type Posts map[PostId]Post
