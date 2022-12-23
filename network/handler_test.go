package network

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"spacemoon/login"
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
	postRequest := httptest.NewRequest(http.MethodPost, "/", http.NoBody)
	postRequest.Header.Set("Content-Type", "multipart/form-data; boundary=*")
	form := url.Values{}
	form.Add("caption", testCaption)
	postRequest.Form = form
	postRecorder := httptest.NewRecorder()
	h.ServeHTTP(postRecorder, postRequest)
	if code := postRecorder.Code; code != http.StatusOK && code != http.StatusAccepted {
		t.Fatalf("invalid status %d", postRecorder.Code)
	}

	getRequest := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	getRecorder := httptest.NewRecorder()
	h.ServeHTTP(getRecorder, getRequest)
	if getRecorder.Code != http.StatusOK {
		t.Fatalf("invalid status %d", postRecorder.Code)
	}

	posts := post.Posts{}
	err := json.Unmarshal(getRecorder.Body.Bytes(), &posts)
	if err != nil {
		t.Fatal(err.Error())
	}

	var newPostId post.Id = ""
	for id, p := range posts {
		if p.GetCaption() == testCaption {
			newPostId = id
		}
	}

	putRequest := httptest.NewRequest(http.MethodPut, "/?id="+string(newPostId), http.NoBody)
	putRecorder := httptest.NewRecorder()
	h.ServeHTTP(putRecorder, putRequest)
	if code := putRecorder.Code; code != http.StatusOK {
		t.Fatalf("unexpected status: %d - %s", code, putRecorder.Body.Bytes())
	}

	var isLiked bool
	err = json.Unmarshal(putRecorder.Body.Bytes(), &isLiked)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !isLiked {
		t.Fatal("post should be liked")
	}

	secondPutRequest := httptest.NewRequest(http.MethodPut, "/like?id="+string(newPostId), http.NoBody)
	secondPutRecorder := httptest.NewRecorder()
	h.ServeHTTP(secondPutRecorder, secondPutRequest)
	if code := secondPutRecorder.Code; code != http.StatusOK {
		t.Fatalf("unexpected status: %d - %s", code, secondPutRecorder.Body.Bytes())
	}
	err = json.Unmarshal(secondPutRecorder.Body.Bytes(), &isLiked)
	if err != nil {
		t.Fatal(err.Error())
	}
	if isLiked {
		t.Fatalf("post should not be liked anymore %s", secondPutRecorder.Body.Bytes())
	}
}

func TestHandler_ServeHTTP_FailsOnPersistenceFail(t *testing.T) {
	h := New(failPersistence{}, stubLoginPersistence{}, stubMediaFilePersistence{})
	request := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	recorder := httptest.NewRecorder()
	h.ServeHTTP(recorder, request)
	if recorder.Code >= 200 && recorder.Code < 300 {
		t.Fatal("error not thrown")
	}
}

func TestHandler_ServeHTTP_PostFailsOnLoginPersistenceFail(t *testing.T) {
	h := New(stubPersistence{}, failLoginPersistence{}, stubMediaFilePersistence{})
	request := httptest.NewRequest(http.MethodPost, "/", http.NoBody)
	recorder := httptest.NewRecorder()
	h.ServeHTTP(recorder, request)
	if recorder.Code >= 200 && recorder.Code < 300 {
		t.Fatalf("error not thrown: %+v", recorder)
	}
}

func TestHandler_ServeHTTP_PutFlipsLikedState(t *testing.T) {

}

type failPersistence struct {
}

func (f failPersistence) AddPost(_ post.Post) error {
	return errors.New("some fake error")
}

func (f failPersistence) GetAllPosts() (post.Posts, error) {
	return nil, errors.New("some fake error")
}

type failLoginPersistence struct {
}

func (f failLoginPersistence) SetUserToken(_ login.UserName, _ login.Token, _ time.Duration) error {
	return fakeError
}

func (f failLoginPersistence) GetUser(_ login.Token) (login.UserName, error) {
	return "", fakeError
}

func (f failLoginPersistence) SignUpUser(u login.UserName, p login.Password) error {
	return fakeError
}

func (f failLoginPersistence) ValidateCredentials(u login.UserName, p login.Password) bool {
	//TODO implement me
	panic("implement me")
}

func (f failLoginPersistence) DeleteUser(name login.UserName) error {
	return fakeError
}

func (f failLoginPersistence) Check(name login.UserName) (bool, error) {
	return false, fakeError
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

func (f stubLoginPersistence) Check(_ login.UserName) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (f stubLoginPersistence) SetUserToken(_ login.UserName, _ login.Token, _ time.Duration) error {
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

func (s stubMediaFilePersistence) SaveFiles(_ map[string]io.Reader, _ string) (post.ContentURIS, error) {
	return nil, nil
}

func (s stubMediaFilePersistence) GetFile(uri string) (io.Reader, error) {
	//TODO implement me
	panic("implement me")
}

func (s stubMediaFilePersistence) Delete(uri string) error {
	//TODO implement me
	panic("implement me")
}

const testAuthor = "Edgar Allan post"
const testCaption = "some caption"

var fakeError = errors.New("some fake error")
