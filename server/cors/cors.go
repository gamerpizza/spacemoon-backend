package cors

import (
	"net/http"
	"strings"
)

func EnableCors(h http.Handler, methods ...string) http.Handler {
	allowed := "OPTIONS, "
	for _, method := range methods {
		allowed += method + ", "
	}
	allowed = strings.TrimRight(allowed, ", ")
	return corsHandler{
		allowed: allowed,
		handler: h,
	}
}

type corsHandler struct {
	allowed string
	handler http.Handler
}

func (c corsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		w.Header().Add("Access-Control-Allow-Methods", c.allowed)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	c.handler.ServeHTTP(w, r)
}
