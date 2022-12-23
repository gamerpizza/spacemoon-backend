package product

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (h *handler) getAllProducts() {
	var products Products
	products, err := h.productPersistence.GetProducts()
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

func (h *handler) getSpecificProduct(productId Id) {
	products, err := h.productPersistence.GetProducts()
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

func (h *handler) respondWithCreatedProduct(createdProduct Product) {
	h.writer.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(h.writer).Encode(createdProduct)
	if err != nil {
		h.writer.WriteHeader(http.StatusInternalServerError)
		_, _ = h.writer.Write([]byte(err.Error()))
	}
}

func (h *handler) getProductFromRequest() (Dto, error) {
	var newProduct = Dto{}
	err := json.NewDecoder(h.request.Body).Decode(&newProduct)
	if err != nil {
		return Dto{}, err
	}
	return newProduct, nil
}

func (h *handler) getIdFromRequest() Id {
	return Id(h.request.FormValue("id"))
}

func isNotEmpty(productId Id) bool {

	return strings.TrimSpace(string(productId)) != ""
}

func isEmpty(productId Id) bool {
	return strings.TrimSpace(string(productId)) == ""
}
