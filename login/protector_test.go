package login

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestProtectedHandler_RequiresToken(t *testing.T) {
	protectedHandler := makeProtectedHandler()

	spy, request, err := prepareRequest(http.MethodGet, "/whatever", http.NoBody)
	if err != nil {
		t.Fatal(err.Error())
	}
	protectedHandler.ServeHTTP(&spy, request)
	badRequestInResponseHeader := checkForStatus(spy, http.StatusBadRequest)
	if !badRequestInResponseHeader {
		t.Fatal("bad request uncaught (no bearer token)")
	}
}

func TestProtected_ServeHTTP_ProxiesTheProtectedHandler(t *testing.T) {
	protectedHandler := makeProtectedHandler()
	spy, request, err := prepareRequest(http.MethodGet, "/whatever", http.NoBody)
	if err != nil {
		t.Fatal(err.Error())
	}
	var bearer = "Bearer <ACCESS TOKEN HERE> "
	request.Header.Add("Authorization", bearer)
	protectedHandler.ServeHTTP(&spy, request)
	if !strings.Contains(spy.body, serveMessage) {
		t.Fatal("protected handler was not called")
	}
}

func makeProtectedHandler() http.Handler {
	p := NewProtector()
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

	_, _ = w.Write([]byte(serveMessage))
}

const serveMessage = "Serving HTTP"
