package product

import (
	"encoding/json"
	"net/http"
	"spacemoon/product"
	"strings"
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
	writer      http.ResponseWriter
	request     *http.Request
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.request = r
	h.writer = w
	switch r.Method {
	case http.MethodPost:
		h.createProduct()
	case http.MethodGet:
		h.getProducts()
	case http.MethodDelete:
		deleteProduct()
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func deleteProduct() {

}

func (h *Handler) getProducts() {
	productId := product.Id(h.request.FormValue("id"))
	if strings.TrimSpace(string(productId)) != "" {
		var p product.Dto = h.persistence.GetProducts()[productId]
		err := json.NewEncoder(h.writer).Encode(p)
		if err != nil {
			h.writer.WriteHeader(http.StatusBadRequest)
			_, _ = h.writer.Write([]byte(err.Error()))
		}
		h.writer.WriteHeader(http.StatusOK)
		return
	}
	var products product.Products = h.persistence.GetProducts()
	err := json.NewEncoder(h.writer).Encode(products)
	if err != nil {
		h.writer.WriteHeader(http.StatusBadRequest)
		_, _ = h.writer.Write([]byte(err.Error()))
	}
	h.writer.WriteHeader(http.StatusOK)
}

func (h *Handler) createProduct() {
	var newProduct = product.Dto{}
	err := json.NewDecoder(h.request.Body).Decode(&newProduct)
	if err != nil {
		h.writer.WriteHeader(http.StatusBadRequest)
		_, _ = h.writer.Write([]byte(err.Error()))
	}
	h.writer.WriteHeader(http.StatusCreated)
}
