package login

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestProtectedHandler_RequiresToken(t *testing.T) {
	protectedHandler := makeProtectedHandler()

	spy, request, err := prepareRequest(http.MethodGet, "/whatever", http.NoBody)
	if err != nil {
		t.Fatal(err.Error())
	}
	protectedHandler.ServeHTTP(&spy, request)
	unauthorizedStatusInResponseHeader := checkForStatus(spy, http.StatusUnauthorized)
	if !unauthorizedStatusInResponseHeader {
		t.Fatalf("%+v", spy)
	}
}

func TestProtected_ServeHTTP_ProxiesTheProtectedHandler(t *testing.T) {
	protectedHandler := makeProtectedHandler()
	spy, request, err := prepareRequest(http.MethodGet, "/whatever", http.NoBody)
	if err != nil {
		t.Fatal(err.Error())
	}
	var bearer = "Bearer " + expectedToken
	request.Header.Add("Authorization", bearer)
	protectedHandler.ServeHTTP(&spy, request)
	if !strings.Contains(spy.body, serveMessage) {
		t.Fatal("protected handler was not called")
	}
}

func TestProtected_UsesTokenFromLogin(t *testing.T) {
	fakePersistence := mockLoginPersistence{}
	token := loginExpectedUser(t, &fakePersistence)
	protectedHandler := validateTokenOnProtectedHandler(t, &fakePersistence, token)
	validateBadTokenOnProtectedHandler(t, protectedHandler)
}

func TestProtector_SetUnprotectedMethod(t *testing.T) {
	var p Protector = NewProtector(&mockLoginPersistence{})
	var testHandler http.Handler = fakeHandler{}
	ph := p.Protect(&testHandler)
	ph.Unprotect(http.MethodGet)

	publicSpy, publicRequest, err := prepareRequest(http.MethodGet, "/whatever", http.NoBody)
	if err != nil {
		t.Fatal(err.Error())
	}
	ph.ServeHTTP(&publicSpy, publicRequest)
	accepted := checkForStatus(publicSpy, http.StatusAccepted)
	if !accepted {
		t.Fatalf("%+v", publicSpy)
	}

	privateSpy, privateRequest, err := prepareRequest(http.MethodPost, "/whatever", http.NoBody)
	if err != nil {
		t.Fatal(err.Error())
	}
	ph.ServeHTTP(&privateSpy, privateRequest)
	unauthorized := checkForStatus(privateSpy, http.StatusUnauthorized)
	if !unauthorized {
		t.Fatalf("%+v", privateSpy)
	}
}

func validateBadTokenOnProtectedHandler(t *testing.T, protectedHandler http.Handler) {
	badBearerSpy, badBearerRequest, err := prepareRequest(http.MethodGet, "/whatever", http.NoBody)
	if err != nil {
		t.Fatalf(err.Error())
	}
	badBearerRequest.Header.Add("Authorization", "bad-token")
	protectedHandler.ServeHTTP(&badBearerSpy, badBearerRequest)
	if badBearerSpy.statusCode != http.StatusForbidden {
		t.Fatalf("did not detect incorrect token: %+v", badBearerSpy)
	}
}

func validateTokenOnProtectedHandler(t *testing.T, fakePersistence Persistence, token string) http.Handler {
	p := NewProtector(fakePersistence)
	var testHandler http.Handler = fakeHandler{}
	protectedHandler := p.Protect(&testHandler)
	spy, req, err := prepareRequest(http.MethodGet, "/whatever", http.NoBody)
	if err != nil {
		t.Fatal(err.Error())
	}
	req.Header.Add("Authorization", "Bearer "+token)
	protectedHandler.ServeHTTP(&spy, req)
	if spy.statusCode != http.StatusAccepted {
		t.Fatalf("unexpected response: %+v", spy)
	}
	return protectedHandler
}

func loginExpectedUser(t *testing.T, fakePersistence Persistence) string {
	loginHandler := NewHandler(fakePersistence, DefaultTokenDuration)
	spy, req, err := prepareRequest(http.MethodGet, "/login", http.NoBody)
	if err != nil {
		t.Fatal(err.Error())
	}
	req.SetBasicAuth(expectedUser, expectedPass)
	loginHandler.ServeHTTP(&spy, req)
	if spy.statusCode != 200 {
		t.Fatalf("unexpected response: %+v", spy)
	}
	token := strings.TrimLeft(spy.body, "token: ")
	return token
}

type stubTokenLoginPersistence struct {
}

func (s stubTokenLoginPersistence) SignUpUser(u UserName, p Password) error {
	//TODO implement me
	panic("implement me")
}

func (s stubTokenLoginPersistence) ValidateCredentials(u UserName, p Password) bool {
	//TODO implement me
	panic("implement me")
}

func (s stubTokenLoginPersistence) DeleteUser(name UserName) error {
	//TODO implement me
	panic("implement me")
}

func (s stubTokenLoginPersistence) SetUserToken(user UserName, token Token, expirationTime time.Duration) error {
	//TODO implement me
	panic("implement me")
	return nil
}

func (s stubTokenLoginPersistence) GetUser(t Token) (UserName, error) {
	if t == expectedToken {
		return expectedUser, nil
	}
	return "", TokenNotFoundError
}

func makeProtectedHandler() http.Handler {
	p := NewProtector(&stubTokenLoginPersistence{})
	var h http.Handler = fakeHandler{}
	protectedHandler := p.Protect(&h)
	return protectedHandler
}

func checkForStatus(spy spyWriter, status int) bool {
	return spy.statusCode == status
}

func prepareRequest(method string, url string, body io.Reader) (spyWriter, *http.Request, error) {
	spy := spyWriter{}
	request, err := http.NewRequest(method, url, body)
	return spy, request, err
}

type fakeHandler struct {
}

func (f fakeHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusAccepted)
	_, _ = w.Write([]byte(serveMessage))
}

const serveMessage = "Serving HTTP"

type mockLoginPersistence struct {
	credentials Tokens
}

func (f *mockLoginPersistence) SignUpUser(u UserName, p Password) error {
	//TODO implement me
	panic("implement me")
}

func (f *mockLoginPersistence) ValidateCredentials(u UserName, p Password) bool {
	//TODO implement me
	panic("implement me")
}

func (f *mockLoginPersistence) DeleteUser(name UserName) error {
	//TODO implement me
	panic("implement me")
}

func (f *mockLoginPersistence) SetUserToken(u UserName, t Token, d time.Duration) error {
	if f.credentials == nil {
		f.credentials = make(Tokens)
	}
	f.credentials[t] = TokenDetails{
		User:       u,
		Expiration: time.Now().Add(d),
	}
	return nil
}

func (f *mockLoginPersistence) GetUser(t Token) (UserName, error) {
	details, exist := f.credentials[t]
	if !exist {
		return "", TokenNotFoundError
	}
	if details.Expiration.Before(time.Now()) {
		return "", ExpiredTokenError
	}
	return details.User, nil
}

const expectedToken = "a-working-bearer-token"
