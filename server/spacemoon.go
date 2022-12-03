package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"spacemoon/product"
	"spacemoon/product/category"
	"spacemoon/server/category_handler"
	"spacemoon/server/product_handler"
	"time"
)

func main() {
	setupHandlers()
	listenAndServe()
}

func setupHandlers() {
	log.Default().Print("starting spacemoon server 🚀")
	log.Default().Print("registering server handlers...")
	http.Handle("/product", product_handler.MakeHandler(&temporaryProductPersistence{}))
	http.Handle("/category", category_handler.MakeHandler(&temporaryCategoryPersistence{}))
	log.Default().Print("handler registration done, ready for takeoff")
}

func listenAndServe() {
	log.Default().Print(getRandomSpaceQuote())
	log.Default().Printf("listening on port %d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Default().Fatalf("error while performing listen and serve: %s", err.Error())
	}
}

func getRandomSpaceQuote() string {
	rand.Seed(time.Now().Unix())
	quote := rand.Intn(2)
	switch quote {
	case 0:
		return "“The stars don't look bigger, but they do look brighter.” ― Sally Ride"
	case 1:
		return "“I see Earth! It is so beautiful.” ― Yuri Gagarin"
	default:
		return ""
	}
}

type temporaryProductPersistence struct {
	savedProducts product.Products
}

func (t *temporaryProductPersistence) DeleteProduct(id product.Id) error {
	delete(t.savedProducts, id)
	return nil
}

func (t *temporaryProductPersistence) GetProducts() (product.Products, error) {
	return t.savedProducts, nil
}

func (t *temporaryProductPersistence) SaveProduct(p product.Product) error {
	if t.savedProducts == nil {
		t.savedProducts = make(product.Products)
	}
	t.savedProducts[p.GetId()] = p.DTO()
	return nil
}

type temporaryCategoryPersistence struct {
	categories category.Categories
}

func (t *temporaryCategoryPersistence) DeleteCategory(name category.Name) {
	delete(t.categories, name)
}

func (t *temporaryCategoryPersistence) SaveCategory(dto category.DTO) {
	if t.categories == nil {
		t.categories = make(category.Categories)
	}
	t.categories[dto.Name] = dto
}

func (t *temporaryCategoryPersistence) GetCategories() category.Categories {
	return t.categories
}

const port = 1234
