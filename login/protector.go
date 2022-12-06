package login

import (
	"net/http"
	"strings"
)

func NewProtector() Protector {
	return protector{}
}

type Protector interface {
	Protect(*http.Handler) http.Handler
}

type protector struct {
}

func (p protector) Protect(h *http.Handler) http.Handler {
	return protected{handler: h}
}

type protected struct {
	handler *http.Handler
}

func (p protected) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bearer := r.Header.Get("Authorization")
	if strings.TrimSpace(bearer) == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("no bearer token"))
		return
	}
	protectedHandler := *p.handler
	protectedHandler.ServeHTTP(w, r)
}
