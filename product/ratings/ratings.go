package ratings

import "spacemoon/product"

// NewProductRater creates a new ProductRater to rate a set of product.Product
func NewProductRater() ProductRater {
	return &productRater{}
}

type ProductRater interface {
	AddRating(product.Id, Rating)
	GetRating(product.Id) Rating
}

type productRater struct {
	ratings ratings
}

func (pr *productRater) GetRating(p product.Id) Rating {
	return pr.ratings[p].finalRating
}

func (pr *productRater) AddRating(p product.Id, r Rating) {
	if pr.ratings == nil {
		pr.ratings = make(ratings)
	}
	productRating := pr.ratings[p]
	newFinalRating := calculateNewFinalRating(productRating, r)
	productRating.ratingHistory = append(pr.ratings[p].ratingHistory, r)
	productRating.finalRating = Rating(newFinalRating)

	pr.ratings[p] = productRating
}

func calculateNewFinalRating(oldRating rating, newValue Rating) int {
	amountOfRatingsForProduct := len(oldRating.ratingHistory)
	lastProductRating := int(oldRating.finalRating)
	newRating := (int(newValue) + lastProductRating*amountOfRatingsForProduct) / (amountOfRatingsForProduct + 1)
	return newRating
}

type Rating uint

type ratings map[product.Id]rating
type rating struct {
	ratingHistory []Rating
	finalRating   Rating
}
