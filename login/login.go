package login

import (
	"net/http"
	"strings"
)

func SecureHandler(h http.Handler) http.Handler {
	return SecuredHandler{h}
}

type SecuredHandler struct {
	handler http.Handler
}

func (s SecuredHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	user, pass, ok := request.BasicAuth()
	if !ok {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(user) == "" {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	if strings.TrimSpace(pass) == "" {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	s.handler.ServeHTTP(writer, request)
}
