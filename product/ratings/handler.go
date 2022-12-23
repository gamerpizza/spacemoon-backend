package ratings

import (
	"encoding/json"
	"fmt"
	"net/http"
	"spacemoon/product"
	"strconv"
)

func MakeRankingsHandler(p Persistence) http.Handler {
	return &rankingsHandler{rater: NewProductRater(p)}
}

type rankingsHandler struct {
	rater ProductRater
}

func (rh *rankingsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	productId := r.FormValue("id")
	switch r.Method {
	case http.MethodGet:
		rating := rh.rater.GetRating(product.Id(productId))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(fmt.Sprintf("Rating: %d", rating)))
	case http.MethodPost:
		ratingStr := r.FormValue("rating")
		rating, err := strconv.ParseInt(ratingStr, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("could not network rating: " + err.Error()))
		}
		w.WriteHeader(http.StatusOK)
		rh.rater.AddRating(product.Id(productId), Score(rating))
		_ = json.NewEncoder(w).Encode(rh.rater.GetRating(product.Id(productId)))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
