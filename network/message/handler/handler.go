package handler

import "net/http"

func New() http.Handler {
	return handler{}
}

type handler struct {
}

func (h handler) ServeHTTP(_ http.ResponseWriter, _ *http.Request) {

}
