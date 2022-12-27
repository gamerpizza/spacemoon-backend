package network

import (
	"encoding/json"
	"io"
	"net/http"
	"spacemoon/login"
	"spacemoon/network/post"
	"strings"
)

type handler struct {
	persistence          Persistence
	writer               http.ResponseWriter
	request              *http.Request
	loginPersistence     login.Persistence
	mediaFilePersistence MediaFilePersistence
	manager              MediaFileContentAdder
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.writer = w
	h.request = r

	switch r.Method {
	case http.MethodGet:
		h.getPosts()
	case http.MethodPost:
		h.post()
	case http.MethodPut:
		h.toggleLike()
	case http.MethodDelete:
		h.delete()
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

func (h handler) post() {
	token := strings.TrimPrefix(h.request.Header.Get("Authorization"), "Bearer ")
	user, err := h.loginPersistence.GetUser(login.Token(token))
	if err != nil {
		h.writer.WriteHeader(http.StatusUnauthorized)
		_, _ = h.writer.Write([]byte(err.Error()))
		return
	}

	err = h.request.ParseMultipartForm(32 << 20)
	if err != nil && !strings.Contains(err.Error(), "EOF") {
		h.writer.WriteHeader(http.StatusBadRequest)
		_, _ = h.writer.Write([]byte(err.Error()))
		return
	}

	caption := post.Caption(h.request.FormValue("caption"))
	newPost := NewPost(caption, user, nil)

	files := make(map[string]io.Reader)
	if h.request.MultipartForm != nil {
		fileHeaders, exists := h.request.MultipartForm.File["media"]
		if exists {
			for _, header := range fileHeaders {
				file, err := header.Open()
				if err != nil {
					h.writer.WriteHeader(http.StatusBadRequest)
					_, _ = h.writer.Write([]byte(err.Error()))
					return
				}
				files[header.Filename] = file
			}
		}
	}

	err = h.manager.SaveNewPostWithMedia(newPost, files)
	if err != nil {
		h.writer.WriteHeader(http.StatusInternalServerError)
		_, _ = h.writer.Write([]byte(err.Error()))
		return
	}
	h.writer.WriteHeader(http.StatusAccepted)
}

func (h handler) toggleLike() {
	id := post.Id(h.request.FormValue("id"))
	allPosts, err := h.persistence.GetAllPosts()
	if err != nil {
		return
	}
	p, exists := allPosts[id]
	if !exists {
		return
	}

	token := strings.TrimPrefix(h.request.Header.Get("Authorization"), "Bearer ")
	user, err := h.loginPersistence.GetUser(login.Token(token))
	if err != nil {
		h.writer.WriteHeader(http.StatusUnauthorized)
		_, _ = h.writer.Write([]byte(err.Error()))
		return
	}

	_, isLiked := p.Likes[string(user)]
	if newIsLiked := !isLiked; isLiked {
		p.RemoveLike(user)

		h.writer.WriteHeader(http.StatusOK)
		err := json.NewEncoder(h.writer).Encode(newIsLiked)
		if err != nil {
			return
		}
	} else {
		p.AddLike(user)
		h.writer.WriteHeader(http.StatusOK)
		err := json.NewEncoder(h.writer).Encode(newIsLiked)
		if err != nil {
			return
		}
	}
	err = h.persistence.AddPost(p)
	if err != nil {
		return
	}
}

func (h handler) delete() {
	h.persistence.DeletePost()
}

func New(np Persistence, lp login.Persistence, mfp MediaFilePersistence) http.Handler {
	return handler{
		persistence:          np,
		loginPersistence:     lp,
		mediaFilePersistence: mfp,
		manager:              NewMediaContentManager(np, mfp)}
}
