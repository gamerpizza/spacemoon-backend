package login

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecureHandler_UsernameCannotBeEmpty(t *testing.T) {
	var h http.Handler = fakeHandler{}
	sh := SecureHandler(h)
	spy := spyWriter{}
	request := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	request.SetBasicAuth("", "pass")
	sh.ServeHTTP(&spy, request)
	if spy.statusCode != http.StatusUnauthorized {
		t.Fatal("did not catch the missing username")
	}
}

func TestSecureHandler_PasswordCannotBeEmpty(t *testing.T) {
	var h http.Handler = fakeHandler{}
	sh := SecureHandler(h)
	spy := spyWriter{}
	request := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	request.SetBasicAuth("user", "")
	sh.ServeHTTP(&spy, request)
	if spy.statusCode != http.StatusUnauthorized {
		t.Fatal("did not catch the missing password")
	}
}

func TestSecureHandler_AuthMustBeProvided(t *testing.T) {
	var h http.Handler = fakeHandler{}
	sh := SecureHandler(h)
	spy := spyWriter{}
	request := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	sh.ServeHTTP(&spy, request)
	if spy.statusCode != http.StatusBadRequest {
		t.Fatal("did not catch the missing authorization")
	}
}

type fakeHandler struct {
}

func (f fakeHandler) ServeHTTP(_ http.ResponseWriter, _ *http.Request) {
}

type spyWriter struct {
	statusCode int
}

func (s *spyWriter) Header() http.Header {
	//TODO implement me
	panic("implement me")
}

func (s *spyWriter) Write(_ []byte) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (s *spyWriter) WriteHeader(statusCode int) {
	s.statusCode = statusCode
}
