package ratings

import (
	"spacemoon/product"
)

type Persistence interface {
	ReadRating(product.Id) Rating
	SaveRating(product.Id, Rating)
}

// NewProductRater creates a new ProductRater to rate a set of product.Product
func NewProductRater(p Persistence) ProductRater {
	return &productRater{persistence: p}
}

type ProductRater interface {
	AddRating(product.Id, Score)
	GetRating(product.Id) Rating
}

type productRater struct {
	persistence Persistence
}

func (pr *productRater) GetRating(p product.Id) Rating {
	return pr.persistence.ReadRating(p)
}

func (pr *productRater) AddRating(p product.Id, r Score) {
	productRating := pr.persistence.ReadRating(p)
	newFinalRating := calculateNewFinalRating(productRating, r)
	productRating.History = append(pr.persistence.ReadRating(p).History, r)
	productRating.Score = Score(newFinalRating)
	pr.persistence.SaveRating(p, productRating)
}

func calculateNewFinalRating(old Rating, new Score) int {
	amountOfRatingsForProduct := len(old.History)
	lastProductRating := int(old.Score)
	newRating := (int(new) + lastProductRating*amountOfRatingsForProduct) / (amountOfRatingsForProduct + 1)
	return newRating
}

type Ratings map[product.Id]Rating
type Rating struct {
	History []Score
	Score   Score
}

type Score int
