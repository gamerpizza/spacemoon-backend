package login

import (
	"net/http"
	"strings"
)

func NewProtector(p Persistence) Protector {
	return protector{persistence: p}
}

type Protector interface {
	Protect(*http.Handler) http.Handler
}

type protector struct {
	persistence Persistence
}

func (p protector) Protect(h *http.Handler) http.Handler {
	return protected{handler: h, persistence: p.persistence}
}

type protected struct {
	handler     *http.Handler
	persistence Persistence
}

func (p protected) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bearer := r.Header.Get("Authorization")
	if p.validateToken(w, bearer) {
		return
	}
	protectedHandler := *p.handler
	protectedHandler.ServeHTTP(w, r)
}

func (p protected) validateToken(w http.ResponseWriter, bearer string) bool {
	if validateBearerTokenPresence(w, bearer) {
		return true
	}
	token := extractToken(bearer)
	_, err := p.getUserFromToken(w, token)
	if err != nil {
		return true
	}
	return false
}

func (p protected) getUserFromToken(w http.ResponseWriter, token Token) (User, error) {
	u, err := p.persistence.GetUser(token)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(err.Error()))
		return "", err
	}
	return u, nil
}

func validateBearerTokenPresence(w http.ResponseWriter, bearer string) bool {
	if strings.TrimSpace(bearer) == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("no bearer token"))
		return true
	}
	return false
}

func extractToken(bearer string) Token {
	return Token(strings.TrimPrefix(bearer, "Bearer "))
}
