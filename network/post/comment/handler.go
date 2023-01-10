// Package comment manages comments for posts in the social network
package comment

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"spacemoon/login"
	"spacemoon/network/post"
	"spacemoon/server/cors"
	"strings"
)

// NewHandler returns a http.Handler, with CORS and login protection incorporated to handle comments
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
	all, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return true
	}
	username, err := h.loginPersistence.GetUser(login.Token(strings.TrimPrefix("Bearer ", r.Header.Get("Authorization"))))
	*comment = New(string(username), fmt.Sprintf("%s", all))

	return false
}

const postKey = "post"
