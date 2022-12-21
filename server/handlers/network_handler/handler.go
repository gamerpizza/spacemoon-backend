package network_handler

import (
	"encoding/json"
	"io"
	"net/http"
	"spacemoon/login"
	"spacemoon/network"
	"spacemoon/network/post"
	"strings"
)

type handler struct {
	persistence          network.Persistence
	writer               http.ResponseWriter
	request              *http.Request
	loginPersistence     login.Persistence
	mediaFilePersistence network.MediaFilePersistence
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

	err = h.request.ParseMultipartForm(32 << 20)
	if err != nil {
		h.writer.WriteHeader(http.StatusBadRequest)
		_, _ = h.writer.Write([]byte(err.Error()))
		return
	}

	manager := network.NewMediaContentManager(h.persistence, h.mediaFilePersistence)
	caption := post.Caption(h.request.FormValue("caption"))
	newPost := network.NewPost(caption, user, nil)
	fileHeaders := h.request.MultipartForm.File["media"]
	files := make(map[string]io.Reader)
	for _, header := range fileHeaders {
		file, err := header.Open()
		if err != nil {
			h.writer.WriteHeader(http.StatusBadRequest)
			_, _ = h.writer.Write([]byte(err.Error()))
			return
		}
		files[header.Filename] = file
	}
	err = manager.SaveNewPostWithMedia(newPost, files)
	if err != nil {
		h.writer.WriteHeader(http.StatusInternalServerError)
		_, _ = h.writer.Write([]byte(err.Error()))
		return
	}
	h.writer.WriteHeader(http.StatusAccepted)
}

func New(np network.Persistence, lp login.Persistence, mfp network.MediaFilePersistence) http.Handler {
	return handler{persistence: np, loginPersistence: lp, mediaFilePersistence: mfp}
}
