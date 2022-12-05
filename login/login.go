package login

import (
	"net/http"
	"strings"
)

func NewHandler(p Persistence) http.Handler {
	return handler{persistence: p}
}

type handler struct {
	persistence Persistence
}

func (h handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	user, pass, ok := request.BasicAuth()
	if !ok {
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte(emptyAuthMessage))
		return
	}
	if strings.TrimSpace(user) == "" {
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte(emptyUsernameMessage))
		return
	}
	if strings.TrimSpace(pass) == "" {
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte(emptyPasswordMessage))
		return
	}

}

func NewProtector(p Persistence) Protector {
	return protector{persistence: p}
}

type Protector interface {
	SecureHandler(http.Handler) http.Handler
}

type Persistence interface {
	ValidateCredentials(usr User, p Password) bool
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
