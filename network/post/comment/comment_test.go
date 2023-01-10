package comment

import (
	"reflect"
	"spacemoon/network/post"
	"spacemoon/network/profile"
	"testing"
)

func TestComment(t *testing.T) {
	var c Comment = New(author, text)
	var a profile.Id = profile.Id(c.Post.Author)
	if a != author {
		t.Fatal("invalid author")
	}
	var m string = string(c.Post.Caption)
	if m != text {
		t.Fatal("invalid text on comment message")
	}
}

func TestCommentManager_Post(t *testing.T) {
	persistence, cm := createTestManager()
	var comment Comment = New(author, text)
	const p post.Id = "some-post"

	cm.Post(comment).On(p)
	comments, _ := cm.GetCommentsFor(p)

	var persistenceComments, _ = persistence.GetCommentsFor(p)
	validatePostedCommentIsFound(t, persistenceComments)
	if !reflect.DeepEqual(comments, persistenceComments) {
		t.Fatal("the retrieved persistenceComments do not correspond to the saved comments")
	}
}

func validatePostedCommentIsFound(t *testing.T, comments []Comment) {
	found := false
	for _, c := range comments {
		if c.Post.Author == author && c.Post.Caption == text {
			found = true
		}
	}
	if !found {
		t.Fatal("comment not found on persistence")
	}
}

func createTestManager() (Persistence, Manager) {
	var persistence Persistence = &fakePersistence{}
	var cm Manager = NewManager(persistence)
	return persistence, cm
}

type fakePersistence struct {
	comments map[post.Id][]Comment
}

func (f *fakePersistence) SaveComment(id post.Id, comment Comment) error {
	if f.comments == nil {
		f.comments = map[post.Id][]Comment{}
	}
	f.comments[id] = append(f.comments[id], comment)
	return nil
}

func (f *fakePersistence) GetCommentsFor(id post.Id) ([]Comment, error) {
	return f.comments[id], nil
}

const author = "test-author"
const text = "something comment something something"
