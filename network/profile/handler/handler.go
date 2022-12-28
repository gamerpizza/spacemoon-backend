package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"spacemoon/login"
	"spacemoon/network/profile"
	"spacemoon/server/persistence/firestore"
	"strings"
)

func New(p profile.Persistence, lp login.Persistence) http.Handler {
	return &handler{persistence: p, loginPersistence: lp}
}

type handler struct {
	persistence      profile.Persistence
	loginPersistence login.Persistence
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.fetchProfile(w, r)
	case http.MethodPut:
		h.updateProfile(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h handler) updateProfile(w http.ResponseWriter, r *http.Request) {
	pr, done := h.getProfile(w, r)
	if done {
		return
	}
	var newProfile profile.Profile
	err := json.NewDecoder(r.Body).Decode(&newProfile)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	token := login.Token(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "))
	user, err := h.loginPersistence.GetUser(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	if login.UserName(pr.Id) != user {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("you can only update your own profile"))
		return
	}
	if strings.TrimSpace(newProfile.Motto.String()) != "" {
		pr.Motto = newProfile.Motto
	}
	if strings.TrimSpace(newProfile.UserName.String()) != "" {
		pr.UserName = newProfile.UserName
	}
	if strings.TrimSpace(newProfile.Avatar.Url.String()) != "" {
		pr.Avatar.Url = newProfile.Avatar.Url
	}
	err = h.persistence.SaveProfile(pr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	return
}

func (h handler) fetchProfile(w http.ResponseWriter, r *http.Request) {
	pr, done := h.getProfile(w, r)
	if done {
		return
	}
	err := json.NewEncoder(w).Encode(pr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	return
}

func (h handler) getProfile(w http.ResponseWriter, r *http.Request) (p profile.Profile, stop bool) {
	id := profile.Id(strings.TrimSpace(r.FormValue("id")))
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return profile.Profile{}, true
	}
	pr, err := h.persistence.GetProfile(id)
	if err != nil && errors.Is(err, firestore.NotFoundError) {
		check, err := h.loginPersistence.Check(login.UserName(id))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return profile.Profile{}, true
		}
		if !check {
			w.WriteHeader(http.StatusNotFound)
			return profile.Profile{}, true
		}
		err = h.persistence.SaveProfile(profile.New(id, profile.UserName(id), "", ""))
		if err != nil {
			return profile.Profile{}, false
		}
		newProfile, err := h.persistence.GetProfile(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return profile.Profile{}, true
		}
		return newProfile, false
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return profile.Profile{}, true
	}
	return pr, false
}
