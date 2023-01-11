// Package comment manages comments for posts in the social network
package comment

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"spacemoon/login"
	"spacemoon/network"
	"spacemoon/network/post"
	"spacemoon/server/cors"
	"strings"
)

// NewHandler returns a http.Handler, with CORS and login protection incorporated to handle comments
func NewHandler(lp login.Persistence, postPersistence network.Persistence, p Persistence) http.Handler {
	var h http.Handler = handler{loginPersistence: lp, postPersistence: postPersistence, manager: NewManager(p)}
	protected := login.NewProtector(lp).Protect(&h)
	protected.Unprotect(http.MethodGet)
	return cors.EnableCors(protected, http.MethodGet, http.MethodPost)
}

type handler struct {
	manager          Manager
	loginPersistence login.Persistence
	postPersistence  network.Persistence
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
		h.getCommentsForPost(w, postId)
	case http.MethodPost:
		h.createNewCommentOnPost(w, r, postId)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h handler) createNewCommentOnPost(w http.ResponseWriter, r *http.Request, postId post.Id) {
	if h.checkIfPostExistsBeforeAction(w, postId, "post a comment") {
		return
	}
	var comment Comment
	if h.createComment(w, r, &comment) {
		return
	}
	h.manager.Post(comment).On(postId)
	return
}

func (h handler) getCommentsForPost(w http.ResponseWriter, postId post.Id) {
	if h.checkIfPostExistsBeforeAction(w, postId, "read comments from") {
		return
	}
	comments, err := h.manager.GetCommentsFor(postId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("could not retrieve comments: " + err.Error()))
		return
	}
	err = json.NewEncoder(w).Encode(comments)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("could not parse comments into json: " + err.Error()))
		return
	}
	return
}

func (h handler) checkIfPostExistsBeforeAction(w http.ResponseWriter, postId post.Id, action string) bool {
	exists, err := h.postPersistence.CheckIfPostExists(postId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("could not read existing posts: " + err.Error()))
		return true
	}
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(fmt.Sprintf("post not found, you cannot post a comment %s on a non-existing post", action)))
		return true
	}
	return false
}

func (h handler) createComment(w http.ResponseWriter, r *http.Request, comment *Comment) bool {
	all, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return true
	}
	token := login.Token(r.Header["Authorization"][0])
	username, err := h.loginPersistence.GetUser(token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error() + fmt.Sprintf("\n%s :: %+v", token, r.Header)))
		return true
	}
	*comment = New(string(username), fmt.Sprintf("%s", all))

	return false
}

const postKey = "post"
