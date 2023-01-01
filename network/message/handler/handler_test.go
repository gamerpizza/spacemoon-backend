package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	var h http.Handler = New()
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("", "/", http.NoBody)
	h.ServeHTTP(recorder, request)
}
