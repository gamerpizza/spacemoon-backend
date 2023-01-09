package comment

import (
	"encoding/json"
	"net/http"
	"spacemoon/login"
	"spacemoon/network/post"
	"spacemoon/network/profile"
	"spacemoon/server/cors"
	"strings"
)

func NewHandler(lp login.Persistence, p Persistence) http.Handler {
	var h http.Handler = handler{loginPersistence: lp, manager: NewManager(p)}
	protected := login.NewProtector(lp).Protect(&h)
	protected.Unprotect(http.MethodGet)
	return cors.EnableCors(protected, http.MethodGet, http.MethodPost)
}

type handler struct {
	manager          Manager
	loginPersistence login.Persistence
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	postId := post.Id(r.FormValue(postKey))
	if strings.TrimSpace(string(postId)) == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("empty post Id"))
		return
	}
	switch r.Method {
	case http.MethodGet:
		comments, _ := h.manager.GetCommentsFor(postId)
		err := json.NewEncoder(w).Encode(comments)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("could not parse comments into json: " + err.Error()))
			return
		}
	case http.MethodPost:
		var comment Comment
		if h.createComment(w, r, &comment) {
			return
		}
		h.manager.Post(comment).On(postId)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h handler) createComment(w http.ResponseWriter, r *http.Request, comment *Comment) bool {
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("could not parse comment into json: " + err.Error()))
		return true
	}
	username, err := h.loginPersistence.GetUser(login.Token(strings.TrimPrefix("Bearer ", r.Header.Get("Authorization"))))
	comment.Author = profile.Id(username)
	return false
}

const postKey = "post"
