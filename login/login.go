package login

import (
	"net/http"
	"strings"
	"time"
)

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

func NewProtector(p Persistence) Protector {
	return protector{persistence: p}
}

type Protector interface {
	SecureHandler(http.Handler) http.Handler
}

type Persistence interface {
	ValidateCredentials(usr User, p Password) bool
	GetUser(token Token) (User, error)
	SetUserToken(user User, token Token, timeToLive time.Duration)
}

type protector struct {
	persistence Persistence
}

func (p protector) SecureHandler(h http.Handler) http.Handler {
	return SecuredHandler{handler: h, persistence: p.persistence}
}

type SecuredHandler struct {
	handler     http.Handler
	persistence Persistence
}

func (s SecuredHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	user, pass, ok := request.BasicAuth()
	if !ok {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(user) == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(pass) == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	isAuthenticated := s.persistence.ValidateCredentials(User(user), Password(pass))
	if !isAuthenticated {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	s.handler.ServeHTTP(writer, request)
}

type User string
type Password string

const emptyUsernameMessage = "username cannot be empty"
const emptyPasswordMessage = "password cannot be empty"
const emptyAuthMessage = "you need to provide a username and corresponding password"

var DefaultTokenDuration = time.Hour
