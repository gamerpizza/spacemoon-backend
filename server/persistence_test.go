package main

import (
	"spacemoon/login"
	"spacemoon/product"
	"spacemoon/product/category"
	"spacemoon/product/ratings"
	"testing"
)

func TestGetLoginPersistence(t *testing.T) {
	var _ login.Persistence = getLoginPersistence()
}

func TestGetProductPersistence(t *testing.T) {
	var _ product.Persistence = getProductPersistence()
}

func TestGetProductRatingsPersistence(t *testing.T) {
	var _ ratings.Persistence = getProductRatingsPersistence()
}

func TestGetCategoryPersistence(t *testing.T) {
	var _ category.Persistence = getCategoryPersistence()
}
