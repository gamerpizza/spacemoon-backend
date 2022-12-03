package category_handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"spacemoon/product"
	"spacemoon/product/category"
	"strings"
)

type Persistence interface {
	GetCategories() category.Categories
	SaveCategory(category.DTO)
	DeleteCategory(category.Name)
}

func MakeHandler(p Persistence) http.Handler {
	return handler{persistence: p}
}

type handler struct {
	persistence Persistence
	writer      http.ResponseWriter
	request     *http.Request
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.writer = w
	h.request = r
	switch r.Method {
	case http.MethodGet:
		h.getCategories()
	case http.MethodPost:
		h.createCategory()
	case http.MethodPut:
		h.addProductToCategory()
	case http.MethodDelete:
		h.delete()
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h handler) getCategories() {
	name := h.getNameFromRequest()
	var c category.Categories = h.persistence.GetCategories()
	if isNotSet(name) {
		h.respondWithAllCategories(c)
		return
	}
	if h.validateThatACategoryExistsWithThatName(c, name) {
		return
	}
	h.respondWithSpecificCategory(c, name)
}

func (h handler) createCategory() {
	c, done := h.getCategoryFromRequestBody()
	if done {
		return
	}
	h.saveNewCategory(c)
}

func (h handler) addProductToCategory() {
	name := h.getNameFromRequest()
	if h.checkIfNameIsInvalid(name) {
		return
	}
	cat, exists := h.persistence.GetCategories()[category.Name(name)]
	if !exists {
		h.writer.WriteHeader(http.StatusNotFound)
		_, _ = h.writer.Write([]byte("category not found"))
		return
	}
	p, done := h.getProductFromRequestBody()
	if done {
		return
	}
	h.addProductAndSave(cat, p)
}

func (h handler) delete() {
	var name category.Name = category.Name(h.request.FormValue("name"))
	if isNotSet(string(name)) {
		h.writer.WriteHeader(http.StatusBadRequest)
		_, _ = h.writer.Write([]byte("you did not specify the name of the category to be deleted"))
		return
	}
	var productId product.Id = product.Id(h.request.FormValue("product_id"))
	if isSet(string(productId)) {
		h.deleteProductFromCategory(name, productId)
		return
	}
	h.persistence.DeleteCategory(name)
	h.writer.WriteHeader(http.StatusNoContent)
}

func (h handler) addProductAndSave(cat category.DTO, p product.Dto) {
	cat.AddProduct(p)
	h.persistence.SaveCategory(cat)
	h.writer.WriteHeader(http.StatusNoContent)
}

func (h handler) checkIfNameIsInvalid(name string) bool {
	if strings.TrimSpace(name) == "" {
		h.writer.WriteHeader(http.StatusBadRequest)
		_, _ = h.writer.Write([]byte("the name of the category cannot be empty"))
		return true
	}
	return false
}

func (h handler) saveNewCategory(c category.DTO) {
	h.persistence.SaveCategory(c)
	h.writer.WriteHeader(http.StatusCreated)
}

func (h handler) getCategoryFromRequestBody() (category.DTO, bool) {
	var c category.DTO
	err := json.NewDecoder(h.request.Body).Decode(&c)
	if err != nil {
		h.writer.WriteHeader(http.StatusBadRequest)
		_, _ = h.writer.Write([]byte(err.Error()))
		return category.DTO{}, true
	}
	return c, false
}

func (h handler) getProductFromRequestBody() (product.Dto, bool) {
	var p product.Dto
	err := json.NewDecoder(h.request.Body).Decode(&p)
	if err != nil {
		h.writer.WriteHeader(http.StatusBadRequest)
		_, _ = h.writer.Write([]byte(fmt.Sprintf("could not decode (json) the product from the request: %s", err.Error())))
		return product.Dto{}, true
	}
	return p, false
}

func (h handler) respondWithSpecificCategory(c category.Categories, name string) {
	h.writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(h.writer).Encode(c[category.Name(name)])
}

func (h handler) validateThatACategoryExistsWithThatName(c category.Categories, name string) bool {
	if _, exists := c[category.Name(name)]; !exists {
		h.writer.WriteHeader(http.StatusNotFound)
		return true
	}
	return false
}

func (h handler) respondWithAllCategories(c category.Categories) {
	encoder := json.NewEncoder(h.writer)
	h.writer.WriteHeader(http.StatusOK)
	_ = encoder.Encode(c)
}

func isNotSet(name string) bool {
	return strings.TrimSpace(name) == ""
}

func isSet(name string) bool {
	return strings.TrimSpace(name) != ""
}

func (h handler) getNameFromRequest() string {
	name := h.request.FormValue("name")
	name = strings.Trim(name, "\"")
	return name
}

func (h handler) deleteProductFromCategory(name category.Name, productId product.Id) {
	cat, exists := h.persistence.GetCategories()[name]
	if !exists {
		h.writer.WriteHeader(http.StatusNotFound)
		return
	}
	cat.DeleteProduct(productId)
	h.writer.WriteHeader(http.StatusNoContent)
	return
}
