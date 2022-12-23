package ratings

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMakeRatingsHandler(t *testing.T) {
	var h http.Handler = MakeRankingsHandler(&fakePersistence{})

	spy := spyWriter{}
	const fakeProductID = "product-id"
	req := httptest.NewRequest(http.MethodGet, "/product/rating?id="+fakeProductID, http.NoBody)
	h.ServeHTTP(&spy, req)
	if spy.header != http.StatusOK {
		t.Fatalf("unexpected header: %d", spy.header)
	}

	var rating Rating
	_ = json.Unmarshal([]byte(strings.TrimSuffix(spy.written, "{written:")), &rating)
	var expectedRating = Rating{
		History: nil,
		Score:   0,
	}
	if rating.Score != expectedRating.Score {
		t.Fatalf("bad response: %+v", spy)
	}

	postReq := httptest.NewRequest(http.MethodPost, "/product/rating?id="+fakeProductID+"&rating=5", http.NoBody)
	postSpy := spyWriter{}
	h.ServeHTTP(&postSpy, postReq)
	if postSpy.header != http.StatusOK {
		t.Fatalf("unexpected header: %d", postSpy.header)
	}
	_ = json.Unmarshal([]byte(strings.TrimSuffix(postSpy.written, "{written:")), &rating)
	if rating.Score != 5 {
		t.Fatalf("bad response: %+v", postSpy)
	}

	postReq = httptest.NewRequest(http.MethodPost, "/product/rating?id="+fakeProductID+"&rating=1", http.NoBody)
	postSpy = spyWriter{}
	h.ServeHTTP(&postSpy, postReq)
	if postSpy.header != http.StatusOK {
		t.Fatalf("unexpected header: %d", postSpy.header)
	}
	_ = json.Unmarshal([]byte(strings.TrimSuffix(postSpy.written, "{written:")), &rating)
	if rating.Score != 3 {
		t.Fatalf("bad response: %+v", postSpy)
	}
}

type spyWriter struct {
	written string
	header  int
}

func (s *spyWriter) Header() http.Header {
	//TODO implement me
	panic("implement me")
}

func (s *spyWriter) Write(bytes []byte) (int, error) {
	s.written = s.written + fmt.Sprintf("%s", bytes)
	return len(bytes), nil
}

func (s *spyWriter) WriteHeader(h int) {
	s.header = h
}
