package product_handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMakeRatingsHandler(t *testing.T) {
	var h http.Handler = MakeRankingsHandler()

	spy := spyWriter{}
	const fakeProductID = "product-id"
	req := httptest.NewRequest(http.MethodGet, "/product/rating?id="+fakeProductID, http.NoBody)
	h.ServeHTTP(&spy, req)
	if spy.header != http.StatusOK {
		t.Fatalf("unexpected header: %d", spy.header)
	}
	if !strings.Contains(spy.written, "Rating: 0") {
		t.Fatalf("bad response: %+v", spy)
	}

	postReq := httptest.NewRequest(http.MethodPost, "/product/rating?id="+fakeProductID+"&rating=5", http.NoBody)
	postSpy := spyWriter{}
	h.ServeHTTP(&postSpy, postReq)
	if postSpy.header != http.StatusOK {
		t.Fatalf("unexpected header: %d", postSpy.header)
	}
	if !strings.Contains(postSpy.written, "Rating: 5") {
		t.Fatalf("bad response: %+v", postSpy)
	}

	postReq = httptest.NewRequest(http.MethodPost, "/product/rating?id="+fakeProductID+"&rating=1", http.NoBody)
	postSpy = spyWriter{}
	h.ServeHTTP(&postSpy, postReq)
	if postSpy.header != http.StatusOK {
		t.Fatalf("unexpected header: %d", postSpy.header)
	}
	if !strings.Contains(postSpy.written, "Rating: 3") {
		t.Fatalf("bad response: %+v", postSpy)
	}
}
