package login

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHandler_AuthCannotBeEmpty(t *testing.T) {
	h, spy, request := setUpGetRequest(defaultTokenDuration)
	h.ServeHTTP(&spy, request)
	if spy.statusCode != http.StatusBadRequest || !strings.Contains(spy.body, emptyAuthMessage) {
		t.Fatal("did not catch the missing username")
	}
}

func TestHandler_UsernameCannotBeEmpty(t *testing.T) {
	h, spy, request := setUpGetRequest(defaultTokenDuration)

	request.SetBasicAuth("", "pass")
	h.ServeHTTP(&spy, request)
	if spy.statusCode != http.StatusBadRequest || !strings.Contains(spy.body, emptyUsernameMessage) {
		t.Fatal("did not catch the missing username")
	}
}

func TestHandler_PasswordCannotBeEmpty(t *testing.T) {
	h, spy, request := setUpGetRequest(defaultTokenDuration)

	request.SetBasicAuth("user", "")
	h.ServeHTTP(&spy, request)
	if spy.statusCode != http.StatusBadRequest || !strings.Contains(spy.body, emptyPasswordMessage) {
		t.Fatal("did not catch the missing password")
	}
}

func TestHandler_Auth(t *testing.T) {
	h, spy, request := setUpGetRequest(defaultTokenDuration)

	request.SetBasicAuth(expectedUser, expectedPass)
	h.ServeHTTP(&spy, request)
	if spy.statusCode != http.StatusOK {
		t.Fatal("expected credentials not recognized")
	}
	if !strings.Contains(spy.body, "token") {
		t.Fatal("expected token not received")
	}
	token := Token(strings.TrimLeft(spy.body, "token: "))

	if user, err := h.persistence.GetUser(token); err != nil || user != expectedUser {
		t.Fatal("expected user not associated to token")
	}
}

func TestHandler_TokenExpiration(t *testing.T) {

	h, spy, request := setUpGetRequest(2 * time.Second)

	request.SetBasicAuth(expectedUser, expectedPass)
	h.ServeHTTP(&spy, request)
	if spy.statusCode != http.StatusOK {
		t.Fatal("expected credentials not recognized")
	}
	if !strings.Contains(spy.body, "token") {
		t.Fatal("expected token not received")
	}
	token := Token(strings.TrimLeft(spy.body, "token: "))

	if user, err := h.persistence.GetUser(token); err != nil || user != expectedUser {
		t.Fatal("expected user not associated to token")
	}

	time.Sleep(3 * time.Second)
	if user, err := h.persistence.GetUser(token); !errors.Is(err, tokenExpiredError) || user != "" {
		t.Fatal("token did not expire")
	}
}

func setUpGetRequest(tokenDuration time.Duration) (*handler, spyWriter, *http.Request) {
	var h = NewHandler(&mockPersistence{}, tokenDuration)
	spy := spyWriter{}
	request := httptest.NewRequest(http.MethodGet, "/login", http.NoBody)
	return h.(*handler), spy, request
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

type mockPersistence struct {
	tokens Credentials
}

func (s *mockPersistence) SetUserToken(user User, token Token, timeToLive time.Duration) {
	if s.tokens == nil {
		s.tokens = make(Credentials)
	}
	s.tokens[token] = TokenDetails{
		user:       user,
		expiration: time.Now().Add(timeToLive),
	}
}

func (s *mockPersistence) GetUser(token Token) (User, error) {
	tokenData, exists := s.tokens[token]
	if !exists {
		return "", tokenNotFoundError
	}
	if tokenData.expiration.Before(time.Now()) {
		return "", tokenExpiredError
	}
	return tokenData.user, nil
}

func (s *mockPersistence) ValidateCredentials(u User, p Password) bool {
	if u != expectedUser || p != expectedPass {
		return false
	}
	return true
}

const expectedUser = "expected-user"
const expectedPass = "expected-pass"
