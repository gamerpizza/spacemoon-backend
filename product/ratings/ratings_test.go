package ratings

import (
	"spacemoon/product"
	"testing"
)

func TestProductRater_CalculatesTheMeanRate(t *testing.T) {
	var pr ProductRater = &productRater{}
	p, err := product.New("test-product", 1000, "")
	if err != nil {
		t.Fatal(err.Error())
	}

	pr.AddRating(p, 50)
	if r := pr.GetRating(p); r != 50 {
		t.Fatalf("incorrect rating, expected %d, got %d", 50, r)
	}

	pr.AddRating(p, 10)
	if r := pr.GetRating(p); r != 30 {
		t.Fatalf("incorrect rating, expected %d, got %d", 30, r)
	}

	pr.AddRating(p, 30)
	if r := pr.GetRating(p); r != 30 {
		t.Fatalf("incorrect rating, expected %d, got %d", 30, r)
	}

	pr.AddRating(p, 10)
	if r := pr.GetRating(p); r != 25 {
		t.Fatalf("incorrect rating, expected %d, got %d", 25, r)
	}

	pr.AddRating(p, 50)
	if r := pr.GetRating(p); r != 30 {
		t.Fatalf("incorrect rating, expected %d, got %d", 30, r)
	}
}

type ProductRater interface {
	AddRating(product.Product, Rating)
	GetRating(product.Product) Rating
}

type productRater struct {
	ratings ratings
}

func (pr *productRater) GetRating(p product.Product) Rating {
	return pr.ratings[p].finalRating
}

func (pr *productRater) AddRating(p product.Product, r Rating) {
	if pr.ratings == nil {
		pr.ratings = make(ratings)
	}
	productRating := pr.ratings[p]
	amountOfRatingsForProduct := len(productRating.ratingHistory)
	lastProductRating := int(productRating.finalRating)
	newRating := (int(r) + lastProductRating*amountOfRatingsForProduct) / (amountOfRatingsForProduct + 1)
	productRating.ratingHistory = append(pr.ratings[p].ratingHistory, r)
	productRating.finalRating = Rating(newRating)
	pr.ratings[p] = productRating
}

type Rating uint

type ratings map[product.Product]totalRating
type totalRating struct {
	ratingHistory []Rating
	finalRating   Rating
}
