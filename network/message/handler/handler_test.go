package handler

import (
	"net/http"
	"net/http/httptest"
	"spacemoon/login"
	"testing"
	"time"
)

func TestHandler(t *testing.T) {
	var h http.Handler = New(stubLoginPersistence{})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("", "/", http.NoBody)
	h.ServeHTTP(recorder, request)
}

func TestHandler_ServeHTTP_OPTIONS(t *testing.T) {
	var h http.Handler = New(stubLoginPersistence{})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodOptions, "/", http.NoBody)
	h.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected status ok not received: %d", recorder.Code)
	}
}

func TestHandler_ServeHTTP_Get_ShouldRetrieveAllConversationsForLoggedInUser(t *testing.T) {
	var h http.Handler = New(stubLoginPersistence{})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	request.Header.Add("Authorization", "Bearer {token}")
	h.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status ok not received: %d", recorder.Code)
	}

}

func TestHandler_ServeHTTP_Get_ShouldFailWithoutABearerToken(t *testing.T) {
	var h http.Handler = New(stubLoginPersistence{})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	h.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status unauthorized not received: %d", recorder.Code)
	}
}

func TestHandler_ServeHTTP_Get_ShouldFailWithABadBearerToken(t *testing.T) {
	var h http.Handler = New(stubFailLoginPersistence{})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	request.Header.Add("Authorization", "Bearer {token}")
	h.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected status forbidden not received: %d", recorder.Code)
	}
}

type method string

type stubFailLoginPersistence struct {
}

func (s stubFailLoginPersistence) SetUserToken(name login.UserName, token login.Token, duration time.Duration) error {
	//TODO implement me
	panic("implement me")
}

func (s stubFailLoginPersistence) GetUser(token login.Token) (login.UserName, error) {
	return "", login.TokenNotFoundError
}

func (s stubFailLoginPersistence) SignUpUser(u login.UserName, p login.Password) error {
	//TODO implement me
	panic("implement me")
}

func (s stubFailLoginPersistence) ValidateCredentials(u login.UserName, p login.Password) bool {
	//TODO implement me
	panic("implement me")
}

func (s stubFailLoginPersistence) DeleteUser(name login.UserName) error {
	//TODO implement me
	panic("implement me")
}

func (s stubFailLoginPersistence) Check(name login.UserName) (bool, error) {
	//TODO implement me
	panic("implement me")
}

type stubLoginPersistence struct {
}

func (s stubLoginPersistence) SetUserToken(name login.UserName, token login.Token, duration time.Duration) error {
	//TODO implement me
	panic("implement me")
}

func (s stubLoginPersistence) GetUser(token login.Token) (login.UserName, error) {
	expectedUserName := login.UserName("test-user")
	return expectedUserName, nil
}

func (s stubLoginPersistence) SignUpUser(u login.UserName, p login.Password) error {
	//TODO implement me
	panic("implement me")
}

func (s stubLoginPersistence) ValidateCredentials(u login.UserName, p login.Password) bool {
	//TODO implement me
	panic("implement me")
}

func (s stubLoginPersistence) DeleteUser(name login.UserName) error {
	//TODO implement me
	panic("implement me")
}

func (s stubLoginPersistence) Check(name login.UserName) (bool, error) {
	//TODO implement me
	panic("implement me")
}
