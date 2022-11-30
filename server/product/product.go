package product

import (
	"encoding/json"
	"net/http"
	"spacemoon/product"
)

func MakeHandler(p Persistence) Handler {
	return Handler{persistence: p}
}

// Persistence is used, as expected, to write and read, to be able to save information.
type Persistence interface {
	GetProducts() product.Products
}

// Handler handles all the calls to the server's product API
type Handler struct {
	persistence Persistence
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createProduct()
	case http.MethodGet:
		h.getProduct(w, r)
	case http.MethodDelete:
		deleteProduct()
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func deleteProduct() {

}

func (h Handler) getProduct(w http.ResponseWriter, _ *http.Request) {
	var products product.Products = h.persistence.GetProducts()
	err := json.NewEncoder(w).Encode(products)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
	}
}

func createProduct() {

}
