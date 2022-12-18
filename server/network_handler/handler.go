package network_handler

import (
	"encoding/json"
	"net/http"
	"spacemoon/login"
	"spacemoon/network"
	"strings"
)

type handler struct {
	persistence      network.Persistence
	writer           http.ResponseWriter
	request          *http.Request
	loginPersistence login.Persistence
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.writer = w
	h.request = r

	switch r.Method {
	case http.MethodGet:
		h.getPosts()
	case http.MethodPost:
		h.Post()
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func (h handler) getPosts() {
	posts, err := h.persistence.GetAllPosts()
	if err != nil {
		h.writer.WriteHeader(http.StatusInternalServerError)
		_, _ = h.writer.Write([]byte(err.Error()))
		return
	}
	err = json.NewEncoder(h.writer).Encode(posts)
	if err != nil {
		h.writer.WriteHeader(http.StatusInternalServerError)
		_, _ = h.writer.Write([]byte(err.Error()))
	}
}

func (h handler) Post() {
	token := strings.TrimPrefix(h.request.Header.Get("Authorization"), "Bearer ")

	user, err := h.loginPersistence.GetUser(login.Token(token))
	if err != nil {
		h.writer.WriteHeader(http.StatusUnauthorized)
		_, _ = h.writer.Write([]byte(err.Error()))
		return
	}

	requestPost := network.Post{}
	err = json.NewDecoder(h.request.Body).Decode(&requestPost)
	if err != nil {
		h.writer.WriteHeader(http.StatusBadRequest)
		_, _ = h.writer.Write([]byte(err.Error()))
		return
	}

	newPost := network.NewPost(requestPost.Caption, user, requestPost.URLS)
	err = h.persistence.AddPost(newPost)
	if err != nil {
		h.writer.WriteHeader(http.StatusInternalServerError)
		_, _ = h.writer.Write([]byte(err.Error()))
		return
	}
	err = json.NewEncoder(h.writer).Encode(newPost)
	if err != nil {
		h.writer.WriteHeader(http.StatusInternalServerError)
		_, _ = h.writer.Write([]byte(err.Error()))
		return
	}

}

func New(np network.Persistence, lp login.Persistence) http.Handler {
	return handler{persistence: np, loginPersistence: lp}
}
