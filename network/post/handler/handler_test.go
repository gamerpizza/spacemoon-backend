package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"spacemoon/network/post"
	"testing"
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
	h := New(&mockNetworkPersistence{}, stubLoginPersistence{}, stubMediaFilePersistence{})
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
	h := New(failNetworkPersistence{}, stubLoginPersistence{}, stubMediaFilePersistence{})
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

func TestHandler_ServeHTTP_Delete(t *testing.T) {
	networkPersistence := &mockNetworkPersistence{}
	testPost := post.New(testCaption, testAuthor, nil)
	err := networkPersistence.AddPost(testPost)
	if err != nil {
		log.Fatal(err.Error())
	}

	h := New(networkPersistence, stubLoginPersistence{}, stubMediaFilePersistence{})
	request := httptest.NewRequest(http.MethodDelete, string("/?id="+testPost.GetId()), http.NoBody)
	recorder := httptest.NewRecorder()
	h.ServeHTTP(recorder, request)
	if code := recorder.Code; code != http.StatusNoContent {
		t.Fatalf("unexpected response status code: %d", code)
	}
	posts, err := networkPersistence.GetAllPosts()
	if err != nil {
		log.Fatal(err.Error())
	}
	_, exists := posts[testPost.GetId()]
	if exists {
		t.Fatalf("post not erased: %+v", posts)
	}
}

var expectedPosts = post.Posts{"1": post.Post{}, "2": post.Post{}}

const testAuthor = "Edgar Allan post"
const testCaption = "some caption"

var fakeError = errors.New("some fake error")
