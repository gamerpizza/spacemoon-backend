package message

import (
	"encoding/json"
	"net/http"
	"spacemoon/login"
	"spacemoon/network/profile"
	"spacemoon/server/cors"
	"strings"
)

func NewHandler(mp Persistence, lp login.Persistence) http.Handler {
	m := NewMessenger(mp, lp)
	var newHandler http.Handler = handler{loginPersistence: lp, messenger: m}
	return cors.EnableCors(login.NewProtector(lp).Protect(&newHandler), http.MethodGet, http.MethodPost)
}

type handler struct {
	messenger        Messenger
	loginPersistence login.Persistence
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		h.getUserConversations(w, r)
	case http.MethodPost:
		h.sendDirectMessage(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h handler) sendDirectMessage(w http.ResponseWriter, r *http.Request) {
	userName, err := h.getUserName(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
	}
	var m Message
	err = json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	if strings.TrimSpace(string(m.Recipient)) == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("cannot send a message without a `Recipient`"))
		return
	}
	err = h.messenger.Send(m).From(profile.Id(userName)).To(profile.Id(m.Recipient)).Now()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	return
}

func (h handler) getUserConversations(w http.ResponseWriter, r *http.Request) {
	userName, err := h.getUserName(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	conversations := h.messenger.GetConversationsWith(profile.Id(userName))
	err = json.NewEncoder(w).Encode(conversations)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	return
}

func (h handler) getUserName(r *http.Request) (login.UserName, error) {
	auth := r.Header.Get("Authorization")
	token := login.Token(strings.TrimSpace(auth))
	userName, err := h.loginPersistence.GetUser(token)
	return userName, err
}
