// Package comment manages comments for posts
package comment

import (
	"spacemoon/network/post"
	"spacemoon/network/profile"
)

// NewManager creates a new Manager. It is used like this to protect the implementation and only make public
// the interface.
func NewManager(p Persistence) Manager {
	return commentManager{persistence: p}
}

// Manager can create, read and delete comments
type Manager interface {
	// Post creates Commenter object, to allow the users of Manager to use of the `Post(Comment).On(post.Id)` format
	Post(comment Comment) Commenter
	// GetCommentsFor gets the comments for a given post.Post indicated by the post.Id. It returns an array of Comment
	GetCommentsFor(post.Id) ([]Comment, error)
}

// New instantiates a new comment
func New(author, text string) Comment {
	return Comment{Author: profile.Id(author), Message: text}
}

// Comment is a Message from a profile.Id Author that will be attached by a Manager to a post.Post
type Comment struct {
	Author  profile.Id
	Message string
}

// Commenter is a helper interface to use of the `Post(Comment).On(post.Id)` format
type Commenter interface {
	On(post post.Id)
}

// Persistence defines how the persistence for Comments should be implemented
type Persistence interface {
	GetCommentsFor(id post.Id) ([]Comment, error)
	SaveComment(post.Id, Comment) error
}

type commentManager struct {
	persistence Persistence
}

func (c commentManager) GetCommentsFor(id post.Id) ([]Comment, error) {
	return c.persistence.GetCommentsFor(id)
}

func (c commentManager) Post(comment Comment) Commenter {
	return commenter{persistence: c.persistence, comment: comment}
}

type commenter struct {
	persistence Persistence
	comment     Comment
}

func (c commenter) On(p post.Id) {
	c.persistence.SaveComment(p, c.comment)
}
