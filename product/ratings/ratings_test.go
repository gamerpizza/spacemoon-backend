package ratings

import (
	"spacemoon/product"
	"testing"
)

func TestProductRater_CalculatesTheMeanRate(t *testing.T) {
	var pr ProductRater = NewProductRater(&fakePersistence{})
	p, err := product.New("test-product", 1000, "", "")
	if err != nil {
		t.Fatal(err.Error())
	}

	pr.AddRating(p.GetId(), 50)
	if r := pr.GetRating(p.GetId()); r.Score != 50 {
		t.Fatalf("incorrect Rating, expected %d, got %d", 50, r)
	}

	pr.AddRating(p.GetId(), 10)
	if r := pr.GetRating(p.GetId()); r.Score != 30 {
		t.Fatalf("incorrect Rating, expected %d, got %d", 30, r)
	}

	pr.AddRating(p.GetId(), 30)
	if r := pr.GetRating(p.GetId()); r.Score != 30 {
		t.Fatalf("incorrect Rating, expected %d, got %d", 30, r)
	}

	pr.AddRating(p.GetId(), 10)
	if r := pr.GetRating(p.GetId()); r.Score != 25 {
		t.Fatalf("incorrect Rating, expected %d, got %d", 25, r)
	}

	pr.AddRating(p.GetId(), 50)
	if r := pr.GetRating(p.GetId()); r.Score != 30 {
		t.Fatalf("incorrect Rating, expected %d, got %d", 30, r)
	}
}

func TestProductRater_GetRating_ShouldBeZeroIfEmpty(t *testing.T) {
	var pr ProductRater = &productRater{persistence: &fakePersistence{}}
	p, err := product.New("test-product", 1000, "", "")
	if err != nil {
		t.Fatal(err.Error())
	}
	if r := pr.GetRating(p.GetId()); r.Score != 0 {
		t.Fatalf("empty Rating should be 0, was %d", r)
	}
}

type fakePersistence struct {
	ratings Ratings
}

func (f *fakePersistence) ReadRating(id product.Id) Rating {
	return f.ratings[id]
}

func (f *fakePersistence) SaveRating(id product.Id, rating Rating) {
	if f.ratings == nil {
		f.ratings = make(Ratings)
	}
	f.ratings[id] = rating
}
