package main

import (
	"fmt"
	"log"
	"net/http"
	"spacemoon/product"
	"spacemoon/server/product_handler"
)

func main() {
	setupHandlers()
	listenAndServe()
}

func setupHandlers() {
	log.Default().Print("starting spacemoon server 🚀")
	log.Default().Print("registering server handlers...")
	http.Handle("/product", product_handler.MakeHandler(&temporaryProductPersistence{}))
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

const port = 1234
