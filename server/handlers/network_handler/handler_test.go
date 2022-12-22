package network_handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"spacemoon/login"
	"spacemoon/network"
	"spacemoon/network/post"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	var _ http.Handler = New(stubPersistence{}, stubLoginPersistence{}, stubMediaFilePersistence{})
}

func TestHandler_ServeHTTP_GetShouldBeAllowed(t *testing.T) {
	h := New(stubPersistence{}, stubLoginPersistence{}, stubMediaFilePersistence{})
	request := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	recorder := httptest.NewRecorder()
	h.ServeHTTP(recorder, request)
	code := recorder.Code
	if code == http.StatusMethodNotAllowed {
		t.Fatal("GET is not allowed")
	}
	if code != http.StatusOK {
		t.Fatalf("GET did not return a 200 status: %d", code)
	}
}

func TestHandler_ServeHTTP_GetShouldReturnAListOfPosts(t *testing.T) {
	h := New(stubPersistence{}, stubLoginPersistence{}, stubMediaFilePersistence{})
	request := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	recorder := httptest.NewRecorder()
	h.ServeHTTP(recorder, request)

	posts := post.Posts{}
	err := json.Unmarshal(recorder.Body.Bytes(), &posts)
	if err != nil {
		t.Fatal(err.Error())
	}
	for id, _ := range expectedPosts {
		if _, exists := posts[id]; !exists {
			t.Fatal("expected posts not retrieved")
		}
	}
}

func TestHandler_ServeHTTP_PostShouldSaveAPost(t *testing.T) {
	h := New(&mockPersistence{}, stubLoginPersistence{}, stubMediaFilePersistence{})
	const testCaption = "some caption"

	p := network.NewPost(testCaption, testAuthor, nil)
	marshal, err := json.Marshal(p)
	if err != nil {
		t.Fatal(err.Error())
	}

	postRequest := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(marshal))
	postRecorder := httptest.NewRecorder()
	h.ServeHTTP(postRecorder, postRequest)
	if postRecorder.Code != http.StatusOK {
		t.Fatalf("invalid status %d", postRecorder.Code)
	}

	getRequest := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	getRecorder := httptest.NewRecorder()
	h.ServeHTTP(getRecorder, getRequest)
	if getRecorder.Code != http.StatusOK {
		t.Fatalf("invalid status %d", postRecorder.Code)
	}

	posts := post.Posts{}
	err = json.Unmarshal(getRecorder.Body.Bytes(), &posts)
	if err != nil {
		t.Fatal(err.Error())
	}

	found := false
	for _, p := range posts {
		if p.GetCaption() == testCaption && p.GetAuthor() == testAuthor {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("post not found")
	}
}

type stubPersistence struct {
}

func (s stubPersistence) AddPost(_ post.Post) error {
	return nil
}

func (s stubPersistence) GetAllPosts() (post.Posts, error) {
	return expectedPosts, nil
}

var expectedPosts = post.Posts{"1": post.Post{}, "2": post.Post{}}

type mockPersistence struct {
	posts post.Posts
}

func (m *mockPersistence) AddPost(p post.Post) error {
	if m.posts == nil {
		m.posts = make(post.Posts)
	}
	m.posts[p.GetId()] = p
	return nil
}

func (m *mockPersistence) GetAllPosts() (post.Posts, error) {
	return m.posts, nil
}

type stubLoginPersistence struct {
}

func (f stubLoginPersistence) SetUserToken(name login.UserName, token login.Token, duration time.Duration) error {
	//TODO implement me
	panic("implement me")
}

func (f stubLoginPersistence) GetUser(_ login.Token) (login.UserName, error) {
	return testAuthor, nil
}

func (f stubLoginPersistence) SignUpUser(u login.UserName, p login.Password) error {
	//TODO implement me
	panic("implement me")
}

func (f stubLoginPersistence) ValidateCredentials(u login.UserName, p login.Password) bool {
	//TODO implement me
	panic("implement me")
}

func (f stubLoginPersistence) DeleteUser(name login.UserName) error {
	//TODO implement me
	panic("implement me")
}

type stubMediaFilePersistence struct {
}

func (s stubMediaFilePersistence) SaveFiles(files map[string]io.Reader, prefix string) (post.ContentURIS, error) {
	//TODO implement me
	panic("implement me")
}

func (s stubMediaFilePersistence) GetFile(uri string) (io.Reader, error) {
	//TODO implement me
	panic("implement me")
}

func (s stubMediaFilePersistence) Delete(uri string) error {
	//TODO implement me
	panic("implement me")
}

const testAuthor = "Edgar Allan Post"
