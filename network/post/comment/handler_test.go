package comment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"spacemoon/login"
	"spacemoon/network/post"
	"strings"
	"testing"
	"time"
)

func TestHandler_ServeHTTP(t *testing.T) {
	var h http.Handler = NewHandler(stubLoginPersistence{}, &fakePersistence{})
	const postId = "test-post"
	var c = Comment{
		Message: text,
	}
	b, err := json.Marshal(c)
	if err != nil {
		t.Fatal(err.Error())
	}
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/?%s=%s", postKey, postId), bytes.NewReader(b))
	req.Header.Add("Authorization", "Bearer test")
	spy := httptest.NewRecorder()
	h.ServeHTTP(spy, req)
	validateResponseCodeIsOk(t, spy.Code)

	getReq := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/?%s=%s", postKey, postId), bytes.NewReader(b))
	getSpy := httptest.NewRecorder()
	h.ServeHTTP(getSpy, getReq)
	validateResponseCodeIsOk(t, getSpy.Code)

	responseBody := getSpy.Body.Bytes()
	var retrievedComments []Comment
	err = json.Unmarshal(responseBody, &retrievedComments)
	if err != nil {
		t.Fatal(err.Error())
	}

	found := false
	for _, comment := range retrievedComments {
		if comment.Author == author && comment.Message == text {
			found = true
		}
	}
	if !found {
		t.Fatal("comment not found")
	}
}

func TestHandler_ServeHTTP_IsCORSEnabled(t *testing.T) {
	var h http.Handler = NewHandler(stubLoginPersistence{}, &fakePersistence{})
	req := httptest.NewRequest(http.MethodOptions, "/", http.NoBody)
	spy := httptest.NewRecorder()
	h.ServeHTTP(spy, req)
	validateResponseCodeIsOk(t, spy.Code)
	validateCORSOptionsResponse(t, spy)
}

func TestHandler_ServeHTTP_POST_ShouldFailWithoutABearerToken(t *testing.T) {
	var h http.Handler = NewHandler(stubLoginPersistence{}, &fakePersistence{})
	const postId = "test-post"
	var c = Comment{
		Message: text,
	}
	b, err := json.Marshal(c)
	if err != nil {
		t.Fatal(err.Error())
	}
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/?%s=%s", postKey, postId), bytes.NewReader(b))
	spy := httptest.NewRecorder()
	h.ServeHTTP(spy, req)
	if spy.Code != http.StatusUnauthorized {
		t.Fatal("did not detect lack of authorization")
	}
}

func validateCORSOptionsResponse(t *testing.T, spy *httptest.ResponseRecorder) {
	if accessControl := spy.Header().Get("Access-Control-Allow-Origin"); accessControl != "*" {
		t.Fatalf("did not return * on Access-Control-Allow-Origin: %s", accessControl)
	}
	validateMethodIsAllowedOnCORS(t, spy, http.MethodGet)
	validateMethodIsAllowedOnCORS(t, spy, http.MethodPost)
}

func validateResponseCodeIsOk(t *testing.T, responseCode int) {
	if responseCode < 200 || responseCode > 299 {
		t.Fatalf("unexpected response code: %d", responseCode)
	}
}

func validateMethodIsAllowedOnCORS(t *testing.T, spy *httptest.ResponseRecorder, method string) {
	if allowedMethods := spy.Header().Get("Access-Control-Allow-Methods"); !strings.Contains(allowedMethods, method) {
		t.Fatalf("allowed methods do not include %s: %s", method, allowedMethods)
	}
}

type spyManager struct {
	manager Manager
	calls   []string
}

func (s spyManager) Post(comment Comment) Commenter {
	s.calls = append(s.calls, "post")
	return s.manager.Post(comment)
}

func (s spyManager) GetCommentsFor(id post.Id) ([]Comment, error) {
	//TODO implement me
	panic("implement me")
}

type stubLoginPersistence struct {
}

func (f stubLoginPersistence) SetUserToken(name login.UserName, token login.Token, duration time.Duration) error {
	//TODO implement me
	panic("implement me")
}

func (f stubLoginPersistence) GetUser(_ login.Token) (login.UserName, error) {
	return author, nil
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

func (f stubLoginPersistence) Check(name login.UserName) (bool, error) {
	//TODO implement me
	panic("implement me")
}
