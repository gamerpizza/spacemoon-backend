package login

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// NewHandler generates a login handler with a reference to the authorization Persistence and a time.Duration
// that defines how long tokens are expected to live once generated
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
	case http.MethodPost:
		h.signup()
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
	isValidAuth := h.persistence.ValidateCredentials(UserName(user), Password(pass))
	if isValidAuth {
		newToken := NewTokenGenerator().NewToken(UseDefaultSize)

		err := h.persistence.SetUserToken(UserName(user), newToken, h.tokenExpirationTime)
		if err != nil {
			h.writer.WriteHeader(http.StatusInternalServerError)
			_, _ = h.writer.Write([]byte(err.Error()))
			return
		}
		h.respondWithAccessToken(newToken)
		return
	}
	h.writer.WriteHeader(http.StatusUnauthorized)
}

func (h *handler) signup() {
	decoder := json.NewDecoder(h.request.Body)
	var u User
	err := decoder.Decode(&u)
	if err != nil {
		h.writer.WriteHeader(http.StatusBadRequest)
		_, _ = h.writer.Write([]byte(err.Error()))
		return
	}
	var exists, _ = h.persistence.Check(u.UserName)
	if exists {
		h.writer.WriteHeader(http.StatusBadRequest)
		_, _ = h.writer.Write([]byte(fmt.Sprintf("username %s is already in use", u.UserName)))
		return
	}
	err = h.persistence.SignUpUser(u.UserName, u.Password)
	if err != nil {
		h.writer.WriteHeader(http.StatusInternalServerError)
		_, _ = h.writer.Write([]byte(err.Error()))
		return
	}
	h.writer.WriteHeader(http.StatusNoContent)
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
