package main

import (
	"spacemoon/login"
	"spacemoon/product/ratings"
	"spacemoon/server/product_handler"
	"testing"
)

func TestGetLoginPersistence(t *testing.T) {
	var _ login.Persistence = GetLoginPersistence()
}

func TestGetProductPersistence(t *testing.T) {
	var _ product_handler.Persistence = GetProductPersistence()
}

func TestGetProductRatingsPersistence(t *testing.T) {
	var _ ratings.Persistence = GetProductRatingsPersistence()
}
