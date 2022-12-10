package login

import (
	"errors"
	"net/http"
	"strings"
	"time"
)

type Persistence interface {
	SetUserToken(user User, token Token, expirationTime time.Duration)
	GetUser(Token) (User, error)
	SignUpUser(u User, p Password)
	ValidateCredentials(u User, p Password) bool
}

func NewHandler(p Persistence, tokenDuration time.Duration) http.Handler {
	return &handler{persistence: p, tokenExpirationTime: tokenDuration}
}

type handler struct {
	persistence         Persistence
	writer              http.ResponseWriter
	request             *http.Request
	tokenExpirationTime time.Duration
}

func (h *handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	h.writer = writer
	h.request = request

	switch request.Method {
	case http.MethodGet:
		h.login()
	default:
		writer.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *handler) login() {
	user, pass, ok := h.request.BasicAuth()
	if h.validateBasicAuth(ok) {
		return
	}
	if h.validateUserName(user) {
		return
	}
	if h.validatePassword(pass) {
		return
	}
	newToken := NewTokenGenerator().NewToken(UseDefaultSize)

	h.persistence.SetUserToken(User(user), newToken, h.tokenExpirationTime)
	h.respondWithAccessToken(newToken)
}

func (h *handler) respondWithAccessToken(token Token) {
	h.writer.WriteHeader(http.StatusOK)
	_, _ = h.writer.Write([]byte("token: " + token))
}

func (h *handler) validatePassword(pass string) bool {
	if strings.TrimSpace(pass) == "" {
		h.writer.WriteHeader(http.StatusBadRequest)
		_, _ = h.writer.Write([]byte(emptyPasswordMessage))
		return true
	}
	return false
}

func (h *handler) validateUserName(user string) bool {
	if strings.TrimSpace(user) == "" {
		h.writer.WriteHeader(http.StatusBadRequest)
		_, _ = h.writer.Write([]byte(emptyUsernameMessage))
		return true
	}
	return false
}

func (h *handler) validateBasicAuth(ok bool) bool {
	if !ok {
		h.writer.WriteHeader(http.StatusBadRequest)
		_, _ = h.writer.Write([]byte(emptyAuthMessage))
		return true
	}
	return false
}

type User string
type Password string

const emptyUsernameMessage = "username cannot be empty"
const emptyPasswordMessage = "password cannot be empty"
const emptyAuthMessage = "you need to provide a username and corresponding password"

var TokenNotFoundError = errors.New("token not found")
var ExpiredTokenError = errors.New("expired token")

var DefaultTokenDuration = time.Hour

type Credentials map[User]Password
