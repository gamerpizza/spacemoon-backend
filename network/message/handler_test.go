package message

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"spacemoon/login"
	"spacemoon/network/profile"
	"testing"
	"time"
)

func TestHandler(t *testing.T) {
	var h http.Handler = NewHandler(stubMessagePersistence{}, stubLoginPersistence{})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("", "/", http.NoBody)
	h.ServeHTTP(recorder, request)
}

func TestHandler_ServeHTTP_OPTIONS(t *testing.T) {
	var h http.Handler = NewHandler(stubMessagePersistence{}, stubLoginPersistence{})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodOptions, "/", http.NoBody)
	h.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected status ok not received: %d", recorder.Code)
	}
}

func TestHandler_ServeHTTP_Get_ShouldRetrieveAllConversationsForLoggedInUser(t *testing.T) {
	var h http.Handler = NewHandler(stubMessagePersistence{}, stubLoginPersistence{})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	request.Header.Add("Authorization", "Bearer {token}")
	h.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status ok not received: %d", recorder.Code)
	}
	var messages map[profile.Id]Conversation
	err := json.NewDecoder(recorder.Body).Decode(&messages)
	if err != nil {
		t.Fatalf("%s -- %+v -- %s", recorder.Body, messages, err.Error())
	}
	fmt.Printf("%+v", messages)
}

func TestHandler_ServeHTTP_Get_ShouldFailWithoutABearerToken(t *testing.T) {
	var h http.Handler = NewHandler(stubMessagePersistence{}, stubLoginPersistence{})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	h.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status unauthorized not received: %d", recorder.Code)
	}
}

func TestHandler_ServeHTTP_Get_ShouldFailWithABadBearerToken(t *testing.T) {
	var h http.Handler = NewHandler(stubMessagePersistence{}, stubFailLoginPersistence{})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	request.Header.Add("Authorization", "Bearer {token}")
	h.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected status forbidden not received: %d", recorder.Code)
	}
}

func TestHandler_ServeHTTP_POST_ShouldFailWithABadBearerToken(t *testing.T) {
	var h http.Handler = NewHandler(stubMessagePersistence{}, stubFailLoginPersistence{})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/", http.NoBody)
	request.Header.Add("Authorization", "Bearer {token}")
	h.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected status forbidden not received: %d", recorder.Code)
	}
}

func TestHandler_ServeHTTP_POST_ShouldFailWithoutABearerToken(t *testing.T) {
	var h http.Handler = NewHandler(stubMessagePersistence{}, stubLoginPersistence{})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/", http.NoBody)
	h.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status unauthorized not received: %d", recorder.Code)
	}
}

func TestHandler_ServeHTTP_POST_ShouldCreateANewMessage(t *testing.T) {
	p := &fakePersistence{}
	var h http.Handler = NewHandler(p, stubLoginPersistence{})
	recorder := httptest.NewRecorder()
	message := Message{
		Recipient: Recipient(u2),
		Content:   "Hello!",
	}
	marshal, err := json.Marshal(message)
	if err != nil {
		return
	}
	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(marshal))
	request.Header.Add("Authorization", "Bearer {token}")
	h.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status ok not received: %d", recorder.Code)
	}

	messages := p.GetMessagesBy(Author("test-user"))
	found := false
	for _, m := range messages[Recipient(u2)] {
		if m.Author == Author("test-user") && m.Content == "Hello!" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected message not found: %+v", messages)
	}
}

type stubFailLoginPersistence struct {
}

func (s stubFailLoginPersistence) SetUserToken(_ login.UserName, _ login.Token, _ time.Duration) error {
	//TODO implement me
	panic("implement me")
}

func (s stubFailLoginPersistence) GetUser(_ login.Token) (login.UserName, error) {
	return "", login.TokenNotFoundError
}

func (s stubFailLoginPersistence) SignUpUser(_ login.UserName, _ login.Password) error {
	//TODO implement me
	panic("implement me")
}

func (s stubFailLoginPersistence) ValidateCredentials(_ login.UserName, _ login.Password) bool {
	//TODO implement me
	panic("implement me")
}

func (s stubFailLoginPersistence) DeleteUser(_ login.UserName) error {
	//TODO implement me
	panic("implement me")
}

func (s stubFailLoginPersistence) Check(_ login.UserName) (bool, error) {
	//TODO implement me
	panic("implement me")
}

type stubLoginPersistence struct {
}

func (s stubLoginPersistence) SetUserToken(_ login.UserName, _ login.Token, _ time.Duration) error {
	//TODO implement me
	panic("implement me")
}

func (s stubLoginPersistence) GetUser(_ login.Token) (login.UserName, error) {
	expectedUserName := login.UserName("test-user")
	return expectedUserName, nil
}

func (s stubLoginPersistence) SignUpUser(_ login.UserName, _ login.Password) error {
	//TODO implement me
	panic("implement me")
}

func (s stubLoginPersistence) ValidateCredentials(_ login.UserName, p login.Password) bool {
	//TODO implement me
	panic("implement me")
}

func (s stubLoginPersistence) DeleteUser(_ login.UserName) error {
	//TODO implement me
	panic("implement me")
}

func (s stubLoginPersistence) Check(u login.UserName) (bool, error) {

	if u == stubLoginBadReceiver {
		return false, nil
	}
	return true, nil
}

const stubLoginBadReceiver = "bad-receiver"

type stubMessagePersistence struct {
}

func (s stubMessagePersistence) GetMessagesBy(_ Author) SentUserMessages {
	return nil
}

func (s stubMessagePersistence) GetMessagesBetween(_ profile.Id) ConversationGetter {
	//TODO implement me
	panic("implement me")
}

func (s stubMessagePersistence) Save(_ Message) error {
	//TODO implement me
	panic("implement me")
}

func (s stubMessagePersistence) GetMessagesFor(_ Recipient) ReceivedUserMessages {
	return nil
}
