package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"spacemoon/login"
	"spacemoon/network/profile"
	"testing"
	"time"
)

func TestHandler(t *testing.T) {
	var h http.Handler = New(getFakePersistence(), &fakeLoginPersistence{})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/?id="+expectedUserId, http.NoBody)
	h.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("ok expected (200), got %d", recorder.Code)
	}
	var p profile.Profile
	err := json.Unmarshal(recorder.Body.Bytes(), &p)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestHandler_ServeHTTP_GET_ShouldFailWithoutAnId(t *testing.T) {
	var h http.Handler = New(fakePersistence{}, &fakeLoginPersistence{})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	h.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("bad request expected (400), got %d", recorder.Code)
	}
}

func TestHandler_ServeHTTP_GET_ShouldFailIfNotFound(t *testing.T) {
	var h http.Handler = New(fakePersistence{}, &fakeLoginPersistence{})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/?id=nobody", http.NoBody)
	h.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("not found expected (403), got %d", recorder.Code)
	}
}

func TestHandler_ServeHTTP_GET_ShouldCreateAProfileIfUserExistsAndProfileDoesNotExist(t *testing.T) {
	var lp login.Persistence = &fakeLoginPersistence{}
	const testUser = "some-user"
	const testPass = "some pass"
	err := lp.SignUpUser(testUser, testPass)
	if err != nil {
		t.Fatal(err.Error())
	}
	var h http.Handler = New(getFakePersistence(), lp)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/profile?id="+testUser, http.NoBody)
	h.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("ok expected (200), got %d", recorder.Code)
	}
	var p profile.Profile
	err = json.Unmarshal(recorder.Body.Bytes(), &p)
	if err != nil {
		t.Fatal(err.Error())
	}
	if p.Id != testUser {
		t.Fatal("invalid profile id, it should be login the username")
	}
}

func TestHandler_ServeHTTP_PUT(t *testing.T) {
	persistence := getFakePersistence()
	var h http.Handler = New(persistence, &fakeLoginPersistence{})
	recorder := httptest.NewRecorder()
	var profileChanges profile.Profile = profile.Profile{
		Motto:    newMotto,
		UserName: newUserName,
		Avatar:   profile.Avatar{Url: newURL},
	}
	marshal, err := json.Marshal(profileChanges)
	if err != nil {
		t.Fatal(err.Error())
	}
	request := httptest.NewRequest(http.MethodPut, "/?id="+expectedUserId, bytes.NewReader(marshal))
	h.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("ok expected (200), got %d", recorder.Code)
	}

	getProfile, err := persistence.GetProfile(expectedUserId)
	if err != nil {
		t.Fatal(err.Error())
	}
	if getProfile.Motto != newMotto {
		t.Fatal("motto not changed correctly")
	}
	if getProfile.UserName != newUserName {
		t.Fatal("username not changed correctly")
	}
	if getProfile.Avatar.Url != newURL {
		t.Fatal("avatar URL not changed correctly")
	}
}

func TestHandler_ServeHTTP_PUT_ShouldNotChangeId(t *testing.T) {
	persistence := getFakePersistence()
	var h http.Handler = New(persistence, &fakeLoginPersistence{})
	recorder := httptest.NewRecorder()
	var profileChanges profile.Profile = profile.Profile{
		Id: "new Id",
	}
	marshal, err := json.Marshal(profileChanges)
	if err != nil {
		t.Fatal(err.Error())
	}
	request := httptest.NewRequest(http.MethodPut, "/?id="+expectedUserId, bytes.NewReader(marshal))
	h.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("ok expected (200), got %d", recorder.Code)
	}

	getProfile, err := persistence.GetProfile(expectedUserId)
	if err != nil {
		t.Fatal(err.Error())
	}
	if getProfile.Id != expectedUserId {
		t.Fatal("id should not be changed")
	}
}

func TestHandler_ServeHTTP_CannotDELETE(t *testing.T) {
	var h http.Handler = New(fakePersistence{}, &fakeLoginPersistence{})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodDelete, "/", http.NoBody)
	h.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("method not allowed expected (405), got %d", recorder.Code)
	}
}

func getFakePersistence() profile.Persistence {
	return fakePersistence{profiles: map[profile.Id]profile.Profile{expectedUserId: profile.Profile{
		Id:       expectedUserId,
		UserName: expectedUserName,
		Motto:    expectedMotto,
		Avatar:   profile.Avatar{Url: expectedAvatarUrl},
	}}}
}

type fakePersistence struct {
	profiles map[profile.Id]profile.Profile
}

func (s fakePersistence) SaveProfile(p profile.Profile) error {
	if s.profiles == nil {
		s.profiles = make(map[profile.Id]profile.Profile)
	}
	s.profiles[p.Id] = p
	return nil
}

func (s fakePersistence) GetProfile(id profile.Id) (profile.Profile, error) {
	p, exists := s.profiles[id]
	if !exists {
		return profile.Profile{}, NotFoundError
	}
	return p, nil
}

type fakeLoginPersistence struct {
	users map[login.UserName]login.Password
}

func (f *fakeLoginPersistence) SetUserToken(name login.UserName, token login.Token, duration time.Duration) error {
	//TODO implement me
	panic("implement me")
}

func (f *fakeLoginPersistence) GetUser(token login.Token) (login.UserName, error) {
	//TODO implement me
	panic("implement me")
}

func (f *fakeLoginPersistence) SignUpUser(u login.UserName, p login.Password) error {
	if f.users == nil {
		f.users = map[login.UserName]login.Password{}
	}
	f.users[u] = p
	return nil
}

func (f *fakeLoginPersistence) ValidateCredentials(u login.UserName, p login.Password) bool {
	//TODO implement me
	panic("implement me")
}

func (f *fakeLoginPersistence) DeleteUser(name login.UserName) error {
	//TODO implement me
	panic("implement me")
}

func (f *fakeLoginPersistence) Check(name login.UserName) (bool, error) {
	_, exists := f.users[name]
	return exists, nil
}

const expectedUserId = "test-user"
const expectedUserName = "some user"
const expectedMotto = "test everything"
const expectedAvatarUrl = "my-avatar.jpg"
const newMotto = "test-driven development"
const newURL = "another-avatar.jpg"
const newUserName = "same user, different name"
