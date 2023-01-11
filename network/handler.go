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
	loginPersistence     login.Persistence
	mediaFilePersistence MediaFilePersistence
	manager              MediaFileContentAdder
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		postId := strings.TrimSpace(r.FormValue("post"))
		switch {
		case postId != "":
			h.getPost(w, postId)
		default:
			h.getPosts(w)
		}
	case http.MethodPost:
		h.post(w, r)
	case http.MethodPut:
		h.toggleLike(w, r)
	case http.MethodDelete:
		h.delete(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func (h handler) getPosts(writer http.ResponseWriter) {
	posts, err := h.persistence.GetAllPosts()
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte(err.Error()))
		return
	}
	err = json.NewEncoder(writer).Encode(posts)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte(err.Error()))
	}
}

func (h handler) post(writer http.ResponseWriter, request *http.Request) {
	token := strings.TrimPrefix(request.Header.Get("Authorization"), "Bearer ")
	user, err := h.loginPersistence.GetUser(login.Token(token))
	if err != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		_, _ = writer.Write([]byte(err.Error()))
		return
	}

	err = request.ParseMultipartForm(32 << 20)
	if err != nil && !strings.Contains(err.Error(), "EOF") {
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte(err.Error()))
		return
	}

	caption := post.Caption(request.FormValue("caption"))
	newPost := post.New(caption, user, nil)

	files := make(map[string]io.Reader)
	if request.MultipartForm != nil {
		fileHeaders, exists := request.MultipartForm.File["media"]
		if exists {
			for _, header := range fileHeaders {
				file, err := header.Open()
				if err != nil {
					writer.WriteHeader(http.StatusBadRequest)
					_, _ = writer.Write([]byte(err.Error()))
					return
				}
				files[header.Filename] = file
			}
		}
	}

	err = h.manager.SaveNewPostWithMedia(newPost, files)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte(err.Error()))
		return
	}
	writer.WriteHeader(http.StatusAccepted)
}

func (h handler) toggleLike(writer http.ResponseWriter, request *http.Request) {
	id := post.Id(request.FormValue("id"))
	allPosts, err := h.persistence.GetAllPosts()
	if err != nil {
		return
	}
	p, exists := allPosts[id]
	if !exists {
		return
	}

	token := strings.TrimPrefix(request.Header.Get("Authorization"), "Bearer ")
	user, err := h.loginPersistence.GetUser(login.Token(token))
	if err != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		_, _ = writer.Write([]byte(err.Error()))
		return
	}

	_, isLiked := p.Likes[string(user)]
	if newIsLiked := !isLiked; isLiked {
		p.RemoveLike(user)

		writer.WriteHeader(http.StatusOK)
		err := json.NewEncoder(writer).Encode(newIsLiked)
		if err != nil {
			return
		}
	} else {
		p.AddLike(user)
		writer.WriteHeader(http.StatusOK)
		err := json.NewEncoder(writer).Encode(newIsLiked)
		if err != nil {
			return
		}
	}
	err = h.persistence.AddPost(p)
	if err != nil {
		return
	}
}

func (h handler) delete(writer http.ResponseWriter, request *http.Request) {
	id := post.Id(request.FormValue("id"))
	allPosts, err := h.persistence.GetAllPosts()
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte(err.Error()))
		return
	}

	token := strings.TrimPrefix(request.Header.Get("Authorization"), "Bearer ")
	user, err := h.loginPersistence.GetUser(login.Token(token))
	if _, exists := allPosts[id]; !exists {
		writer.WriteHeader(http.StatusNoContent)
		return
	}

	if allPosts[id].Author != user {
		writer.WriteHeader(http.StatusUnauthorized)
		_, _ = writer.Write([]byte("user does not own that post"))
		return
	}

	err = h.persistence.DeletePost(id)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte(err.Error()))
		return
	}
	writer.WriteHeader(http.StatusNoContent)
}

func (h handler) getPost(w http.ResponseWriter, id string) {
	var p post.Post
	var err error
	p, err = h.persistence.GetPost(post.Id(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	err = json.NewEncoder(w).Encode(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
	}
}

func NewHandler(np Persistence, lp login.Persistence, mfp MediaFilePersistence) http.Handler {
	return handler{
		persistence:          np,
		loginPersistence:     lp,
		mediaFilePersistence: mfp,
		manager:              NewMediaContentManager(np, mfp)}
}
