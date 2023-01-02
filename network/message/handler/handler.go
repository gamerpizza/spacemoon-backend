package handler

import (
	"net/http"
	"spacemoon/login"
	"spacemoon/network/message"
	"spacemoon/server/cors"
	"strings"
)

func New(lp login.Persistence) http.Handler {
	var newHandler http.Handler = handler{loginPersistence: lp}

	return cors.EnableCors(login.NewProtector(lp).Protect(&newHandler), http.MethodGet)
}

type handler struct {
	messenger        message.Messenger
	loginPersistence login.Persistence
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		auth := r.Header.Get("Authorization")
		if strings.TrimSpace(auth) == "" {
			w.WriteHeader(http.StatusUnauthorized)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
