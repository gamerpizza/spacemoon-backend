package product_handler

import (
	"encoding/json"
	"net/http"
	"spacemoon/product"
	"strings"
)

func (h *handler) getAllProducts() {
	var products product.Products
	products, err := h.persistence.GetProducts()
	if err != nil {
		h.writer.WriteHeader(http.StatusInternalServerError)
		_, _ = h.writer.Write([]byte(err.Error()))
		return
	}
	h.writer.WriteHeader(http.StatusOK)
	err = json.NewEncoder(h.writer).Encode(products)
	if err != nil {
		h.writer.WriteHeader(http.StatusBadRequest)
		_, _ = h.writer.Write([]byte(err.Error()))
	}
}

func (h *handler) getSpecificProduct(productId product.Id) {
	products, err := h.persistence.GetProducts()
	if err != nil {
		h.writer.WriteHeader(http.StatusInternalServerError)
		_, _ = h.writer.Write([]byte(err.Error()))
	}
	p, exists := products[productId]
	if !exists {
		h.writer.WriteHeader(http.StatusNotFound)
		return
	}
	h.writer.WriteHeader(http.StatusOK)
	err = json.NewEncoder(h.writer).Encode(p)
	if err != nil {
		h.writer.WriteHeader(http.StatusBadRequest)
		_, _ = h.writer.Write([]byte(err.Error()))
	}
}

func (h *handler) respondWithCreatedProduct(createdProduct product.Product) {
	h.writer.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(h.writer).Encode(createdProduct)
	if err != nil {
		h.writer.WriteHeader(http.StatusInternalServerError)
		_, _ = h.writer.Write([]byte(err.Error()))
	}
}

func (h *handler) getProductFromRequest() (product.Dto, error) {
	var newProduct = product.Dto{}
	err := json.NewDecoder(h.request.Body).Decode(&newProduct)
	if err != nil {
		return product.Dto{}, err
	}
	return newProduct, nil
}

func (h *handler) getIdFromRequest() product.Id {
	return product.Id(h.request.FormValue("id"))
}

func isNotEmpty(productId product.Id) bool {

	return strings.TrimSpace(string(productId)) != ""
}

func isEmpty(productId product.Id) bool {
	return strings.TrimSpace(string(productId)) == ""
}
