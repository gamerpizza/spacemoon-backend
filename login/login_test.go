package login

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_AuthCannotBeEmpty(t *testing.T) {
	var h http.Handler = handler{}
	spy := spyWriter{}
	request := httptest.NewRequest(http.MethodGet, "/login", http.NoBody)
	h.ServeHTTP(&spy, request)
	if spy.statusCode != http.StatusBadRequest || !strings.Contains(spy.body, emptyAuthMessage) {
		t.Fatal("did not catch the missing username")
	}
}

func TestHandler_UsernameCannotBeEmpty(t *testing.T) {
	var h http.Handler = handler{}
	spy := spyWriter{}
	request := httptest.NewRequest(http.MethodGet, "/login", http.NoBody)
	request.SetBasicAuth("", "pass")
	h.ServeHTTP(&spy, request)
	if spy.statusCode != http.StatusBadRequest || !strings.Contains(spy.body, emptyUsernameMessage) {
		t.Fatal("did not catch the missing username")
	}
}

func TestHandler_PasswordCannotBeEmpty(t *testing.T) {
	var h http.Handler = handler{}
	spy := spyWriter{}
	request := httptest.NewRequest(http.MethodGet, "/login", http.NoBody)
	request.SetBasicAuth("user", "")
	h.ServeHTTP(&spy, request)
	if spy.statusCode != http.StatusBadRequest || !strings.Contains(spy.body, emptyPasswordMessage) {
		t.Fatal("did not catch the missing password")
	}
}

func TestHandler_Auth(t *testing.T) {
	var testPersistence Persistence = stubPersistence{}
	var h http.Handler = NewHandler(testPersistence)
	spy := spyWriter{}
	request := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	request.SetBasicAuth(expectedUser, expectedPass)
	h.ServeHTTP(&spy, request)
	if spy.statusCode == http.StatusUnauthorized {
		t.Fatal("expected credentials not recognized")
	}
	if !strings.Contains(spy.body, "token") {
		t.Fatal("expected token not received")
	}
}

func TestProtector_Auth(t *testing.T) {
	sh := createSecureHandler()
	testValidCredentials(t, sh)
	testInvalidCredentials(t, sh)
}

func createSecureHandler() http.Handler {
	var pe Persistence = stubPersistence{}
	var pr Protector = NewProtector(pe)
	var h http.Handler = fakeHandler{}
	sh := pr.SecureHandler(h)
	return sh
}

func testValidCredentials(t *testing.T, sh http.Handler) {
	spy := spyWriter{}
	request := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	request.SetBasicAuth(expectedUser, expectedPass)
	sh.ServeHTTP(&spy, request)
	if spy.statusCode == http.StatusUnauthorized {
		t.Fatal("expected credentials not recognized")
	}
}

func testInvalidCredentials(t *testing.T, sh http.Handler) {
	failingRequest := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	failingRequest.SetBasicAuth("bad-user", "bad-pass")
	failingSpy := spyWriter{}
	sh.ServeHTTP(&failingSpy, failingRequest)
	if failingSpy.statusCode != http.StatusUnauthorized {
		t.Fatal("bad credentials not recognized")
	}
}

type fakeHandler struct {
}

func (f fakeHandler) ServeHTTP(_ http.ResponseWriter, _ *http.Request) {
}

type spyWriter struct {
	statusCode int
	body       string
}

func (s *spyWriter) Header() http.Header {
	//TODO implement me
	panic("implement me")
}

func (s *spyWriter) Write(w []byte) (int, error) {
	s.body += fmt.Sprintf("%s%s", s.body, w)
	return len(w), nil
}

func (s *spyWriter) WriteHeader(statusCode int) {
	s.statusCode = statusCode
}

type stubPersistence struct {
}

func (s stubPersistence) ValidateCredentials(u User, p Password) bool {
	if u != expectedUser || p != expectedPass {
		return false
	}
	return true
}

const expectedUser = "expected-user"
const expectedPass = "expected-pass"
