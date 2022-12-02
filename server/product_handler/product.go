// Package product_handler handles calls related to the product.
package product_handler

import (
	"net/http"
	"spacemoon/product"
)

// MakeHandler creates a product handler and attributes it a Persistence. This is made to allow the Persistence
// implementation to be easily changed dynamically.
func MakeHandler(p Persistence) http.Handler {
	return &handler{persistence: p}
}

// Persistence is used, as expected, to write and read, to be able to save information.
type Persistence interface {
	GetProducts() (product.Products, error)
	SaveProduct(product.Product) error
	DeleteProduct(product.Id) error
}

// handler handles all the calls to the server's product API
type handler struct {
	persistence Persistence
	writer      http.ResponseWriter
	request     *http.Request
}

// ServeHTTP will handle the request according to the http method. For http.MethodPost, it will create a new product.
// For http.MethodGet it will retrieve one or more products. For http.MethodPut it will update a product.
// For http.MethodDelete, it will delete a product. All other methods will return a http.StatusMethodNotAllowed header.
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.request = r
	h.writer = w
	switch r.Method {
	case http.MethodPost:
		h.createProduct()
	case http.MethodGet:
		h.getProducts()
	case http.MethodDelete:
		h.deleteProduct()
	case http.MethodPut:
		updateProduct()
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *handler) getProducts() {
	if productId := h.getIdFromRequest(); isNotEmpty(productId) {
		h.getSpecificProduct(productId)
		return
	}

	h.getAllProducts()
}

func (h *handler) createProduct() {
	newProduct, err := h.getProductFromRequest()
	if err != nil {
		h.writer.WriteHeader(http.StatusBadRequest)
		_, _ = h.writer.Write([]byte(err.Error()))
		return
	}
	existingProducts, err := h.persistence.GetProducts()
	if err != nil {
		h.writer.WriteHeader(http.StatusInternalServerError)
		_, _ = h.writer.Write([]byte(err.Error()))
		return
	}
	_, exists := existingProducts[newProduct.Id]
	if exists {
		h.writer.WriteHeader(http.StatusConflict)
		_, _ = h.writer.Write([]byte("a product with that Id already exists"))
		return
	}
	createdProduct, err := product.New(newProduct.Name, newProduct.Price, newProduct.Description)
	if err != nil {
		h.writer.WriteHeader(http.StatusBadRequest)
		_, _ = h.writer.Write([]byte(err.Error()))
		return
	}
	err = h.persistence.SaveProduct(createdProduct)
	if err != nil {
		h.writer.WriteHeader(http.StatusInternalServerError)
		_, _ = h.writer.Write([]byte(err.Error()))
		return
	}
	h.respondWithCreatedProduct(createdProduct)
}

// TODO
func updateProduct() {
}

func (h *handler) deleteProduct() {
	productId := h.getIdFromRequest()
	if isEmpty(productId) {
		h.writer.WriteHeader(http.StatusBadRequest)
		_, _ = h.writer.Write([]byte("you did not specify the Id of the item to be deleted"))
		return
	}
	existingProducts, err := h.persistence.GetProducts()
	if err != nil {
		h.writer.WriteHeader(http.StatusInternalServerError)
		_, _ = h.writer.Write([]byte(err.Error()))
		return
	}
	_, exists := existingProducts[productId]
	if !exists {
		h.writer.WriteHeader(http.StatusConflict)
		_, _ = h.writer.Write([]byte("product not found"))
		return
	}
	err = h.persistence.DeleteProduct(productId)
	if err != nil {
		h.writer.WriteHeader(http.StatusInternalServerError)
		_, _ = h.writer.Write([]byte(err.Error()))
		return
	}
	h.writer.WriteHeader(http.StatusOK)
}
