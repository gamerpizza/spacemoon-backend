package network_handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"spacemoon/network"
	"testing"
)

func TestNew(t *testing.T) {
	var _ http.Handler = New(stubPersistence{})
}

func TestHandler_ServeHTTP_GetShouldBeAllowed(t *testing.T) {
	h := New(stubPersistence{})
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
	h := New(stubPersistence{})
	request := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	recorder := httptest.NewRecorder()
	h.ServeHTTP(recorder, request)

	posts := network.Posts{}
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
	h := New(&mockPersistence{})
	const testCaption = "some caption"
	const testAuthor = "Edgar Allan Post"
	post := network.NewPost(testCaption, testAuthor, nil)
	marshal, err := json.Marshal(post)
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

	posts := network.Posts{}
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

func (s stubPersistence) AddPost(_ network.Post) error {
	return nil
}

func (s stubPersistence) GetAllPosts() (network.Posts, error) {
	return expectedPosts, nil
}

var expectedPosts = network.Posts{"1": network.Post{}, "2": network.Post{}}

type mockPersistence struct {
	posts network.Posts
}

func (m *mockPersistence) AddPost(post network.Post) error {
	if m.posts == nil {
		m.posts = make(network.Posts)
	}
	m.posts[post.GetId()] = post
	return nil
}

func (m *mockPersistence) GetAllPosts() (network.Posts, error) {
	return m.posts, nil
}
