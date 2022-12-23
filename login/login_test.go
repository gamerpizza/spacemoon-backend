package login

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHandler_AuthCannotBeEmpty(t *testing.T) {
	h, spy, request := setUpGetRequest(DefaultTokenDuration)
	h.ServeHTTP(&spy, request)
	if spy.statusCode != http.StatusBadRequest || !strings.Contains(spy.body, emptyAuthMessage) {
		t.Fatal("did not catch the missing username")
	}
}

func TestHandler_UsernameCannotBeEmpty(t *testing.T) {
	h, spy, request := setUpGetRequest(DefaultTokenDuration)

	request.SetBasicAuth("", "pass")
	h.ServeHTTP(&spy, request)
	if spy.statusCode != http.StatusBadRequest || !strings.Contains(spy.body, emptyUsernameMessage) {
		t.Fatal("did not catch the missing username")
	}
}

func TestHandler_PasswordCannotBeEmpty(t *testing.T) {
	h, spy, request := setUpGetRequest(DefaultTokenDuration)

	request.SetBasicAuth("user", "")
	h.ServeHTTP(&spy, request)
	if spy.statusCode != http.StatusBadRequest || !strings.Contains(spy.body, emptyPasswordMessage) {
		t.Fatal("did not catch the missing password")
	}
}

func TestHandler_Auth(t *testing.T) {
	h, spy, request := setUpGetRequest(DefaultTokenDuration)
	token := getTokenFromHTTPCall(t, request, h, spy)
	checkIfTokenIsAssociatedWithExpectedUser(t, h, token)
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

func TestHandler_ServeHTTP_PostShouldSignup(t *testing.T) {
	var h = NewHandler(&mockPersistence{}, time.Hour)
	signUpSpy := spyWriter{}
	const newUsername = "newUser"
	const newPassword = "newPass"
	user := User{
		UserName: "newUser",
		Password: "newPass",
	}
	marshal, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err.Error())
	}

	signUpRequest := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(marshal))
	h.ServeHTTP(&signUpSpy, signUpRequest)

	if signUpSpy.statusCode < 200 || signUpSpy.statusCode >= 300 {
		t.Fatalf("bad status code :%d", signUpSpy.statusCode)
	}

	doubleSignUpSpy := spyWriter{}
	h.ServeHTTP(&doubleSignUpSpy, signUpRequest)
	if doubleSignUpSpy.statusCode >= 200 && doubleSignUpSpy.statusCode < 300 {
		t.Fatalf("bad status code, should not be able to sign up with the same name :%d", doubleSignUpSpy.statusCode)
	}

	loginSpy := spyWriter{}
	loginRequest := httptest.NewRequest(http.MethodGet, "/login", http.NoBody)
	loginRequest.SetBasicAuth(newUsername, newPassword)
	h.ServeHTTP(&loginSpy, loginRequest)

	if loginSpy.statusCode < 200 || loginSpy.statusCode >= 300 {
		t.Fatalf("bad status code :%d", loginSpy.statusCode)
	}

	badLoginSpy := spyWriter{}
	badLoginRequest := httptest.NewRequest(http.MethodGet, "/login", http.NoBody)
	badLoginRequest.SetBasicAuth("bad-user", "pass")
	h.ServeHTTP(&badLoginSpy, badLoginRequest)

	if badLoginSpy.statusCode >= 200 && badLoginSpy.statusCode < 300 {
		t.Fatalf("bad status code, user should not be able to login :%d", badLoginSpy.statusCode)
	}
}

func checkIfTokenIsAssociatedWithExpectedUser(t *testing.T, h *handler, token Token) {
	if user, err := h.persistence.GetUser(token); err != nil || user != expectedUser {
		t.Fatal("expected user not associated to token")
	}
}

func getTokenFromHTTPCall(t *testing.T, request *http.Request, h *handler, spy spyWriter) Token {
	request.SetBasicAuth(expectedUser, expectedPass)
	h.ServeHTTP(&spy, request)
	if spy.statusCode != http.StatusOK {
		t.Fatal("expected credentials not recognized")
	}
	if !strings.Contains(spy.body, "token") {
		t.Fatal("expected token not received")
	}
	token := Token(strings.TrimLeft(spy.body, "token: "))
	return token
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
	tokens      Tokens
	credentials map[UserName]Password
}

func (s *mockPersistence) Check(_ UserName) (bool, error) {
	return false, nil
}

func (s *mockPersistence) SignUpUser(u UserName, p Password) error {
	if s.credentials == nil {
		s.credentials = make(map[UserName]Password)
	}
	s.credentials[u] = p
	return nil
}

func (s *mockPersistence) DeleteUser(name UserName) error {
	delete(s.credentials, name)
	return nil
}

func (s *mockPersistence) SetUserToken(user UserName, token Token, duration time.Duration) error {
	if s.tokens == nil {
		s.tokens = make(Tokens)
	}
	s.tokens[token] = TokenDetails{
		User:       user,
		Expiration: time.Now().Add(duration),
	}
	return nil
}

func (s *mockPersistence) GetUser(token Token) (UserName, error) {
	tokenData, exists := s.tokens[token]
	if !exists {
		return "", TokenNotFoundError
	}
	if tokenData.Expiration.Before(time.Now()) {
		return "", tokenExpiredError
	}
	return tokenData.User, nil
}

func (s *mockPersistence) ValidateCredentials(u UserName, p Password) bool {
	if u == expectedUser || p == expectedPass {
		return true
	}
	if s.credentials != nil && s.credentials[u] == p {
		return true
	}
	return false
}

const expectedUser = "expected-user"
const expectedPass = "expected-pass"
