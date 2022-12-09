package ratings

import (
	"spacemoon/product"
	"testing"
)

func TestProductRater_CalculatesTheMeanRate(t *testing.T) {
	var pr ProductRater = NewProductRater()
	p, err := product.New("test-product", 1000, "")
	if err != nil {
		t.Fatal(err.Error())
	}

	pr.AddRating(p.GetId(), 50)
	if r := pr.GetRating(p.GetId()); r != 50 {
		t.Fatalf("incorrect rating, expected %d, got %d", 50, r)
	}

	pr.AddRating(p.GetId(), 10)
	if r := pr.GetRating(p.GetId()); r != 30 {
		t.Fatalf("incorrect rating, expected %d, got %d", 30, r)
	}

	pr.AddRating(p.GetId(), 30)
	if r := pr.GetRating(p.GetId()); r != 30 {
		t.Fatalf("incorrect rating, expected %d, got %d", 30, r)
	}

	pr.AddRating(p.GetId(), 10)
	if r := pr.GetRating(p.GetId()); r != 25 {
		t.Fatalf("incorrect rating, expected %d, got %d", 25, r)
	}

	pr.AddRating(p.GetId(), 50)
	if r := pr.GetRating(p.GetId()); r != 30 {
		t.Fatalf("incorrect rating, expected %d, got %d", 30, r)
	}
}

func TestProductRater_GetRating_ShouldBeZeroIfEmpty(t *testing.T) {
	var pr ProductRater = &productRater{}
	p, err := product.New("test-product", 1000, "")
	if err != nil {
		t.Fatal(err.Error())
	}
	if r := pr.GetRating(p.GetId()); r != 0 {
		t.Fatalf("empty rating should be 0, was %d", r)
	}
}
