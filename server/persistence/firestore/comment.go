package firestore

import (
	"fmt"
	"spacemoon/network/post"
	"spacemoon/network/post/comment"
)

func (p *fireStorePersistence) GetCommentsFor(id post.Id) ([]comment.Comment, error) {
	collection := p.storage.Collection(commentsCollection).Doc(string(id)).Collection(commentsSubCollection)
	documents, err := collection.Documents(p.ctx).GetAll()
	if err != nil {
		return nil, err
	}
	var comments []comment.Comment
	for _, document := range documents {
		var cmt comment.Comment
		err = document.DataTo(&cmt)
		if err != nil {
			return nil, fmt.Errorf("could not parse saved cmt: %w", err)
		}
		comments = append(comments, cmt)
	}

	return comments, nil
}

func (p *fireStorePersistence) SaveComment(id post.Id, c comment.Comment) error {
	collection := p.storage.Collection(commentsCollection).Doc(string(id)).Collection(commentsSubCollection)
	_, _, err := collection.Add(p.ctx, c)
	if err != nil {
		return fmt.Errorf("could not save comment: %w", err)
	}
	return nil
}
