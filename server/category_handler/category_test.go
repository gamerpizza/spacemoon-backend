package category_handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"spacemoon/product"
	"spacemoon/product/category"
	"strings"
	"testing"
)

func TestCategory_POST(t *testing.T) {
	h := MakeHandler(&fakePersistence{})
	spy := spyWriter{}
	testCategory, expectedCategoryJSON := createTestCategory(t)

	request := makePOSTRequest(t, expectedCategoryJSON)
	h.ServeHTTP(&spy, request)

	validateSavedCategory(t, h, spy, testCategory, expectedCategoryJSON)
}

func TestCategory_PUT(t *testing.T) {
	h, testCategory, _ := postNewCategory(t)
	testProductJSON, _ := addProductToCategory(t, testCategory, h)
	checkThatTheProductIsFoundOnTheCategory(t, h, testCategory, testProductJSON)
}

func TestCategory_GET(t *testing.T) {
	h := makeStubPersistenceHandler()
	spy := spyWriter{}
	h.ServeHTTP(&spy, httptest.NewRequest(http.MethodGet, "/category", http.NoBody))

	if fmt.Sprintf("%s", spy.written) != fmt.Sprintf("%s", expectedCategoriesJSON) {
		t.Fatalf("did not retrieve expected categories, \nexpected '%s'\nreceived '%s'\n", expectedCategoriesJSON, spy.written)
	}
}

func TestCategory_GETByName(t *testing.T) {
	const expectedCategoryName = "hot products"
	h := makeStubPersistenceHandler()
	spy := spyWriter{}
	performRequestToGetCategoryByName(t, h, &spy, expectedCategoryName)
	validateExpectedCategoryFromStubPersistenceInResponse(t, &spy, expectedCategoryName)
}

func TestCategory_DELETE(t *testing.T) {
	h, testCategory, expectedCategoryJSON := postNewCategory(t)
	spy := spyWriter{}
	validateSavedCategory(t, h, spy, testCategory, expectedCategoryJSON)
	deleteCategory(t, h, testCategory.Name)
	validateCategoryDeletion(t, h, testCategory, expectedCategoryJSON)
}

func TestCategory_DELETE_PRODUCT(t *testing.T) {
	h, testCategory, _ := postNewCategory(t)
	testProductJSON, productId := addProductToCategory(t, testCategory, h)
	checkThatTheProductIsFoundOnTheCategory(t, h, testCategory, testProductJSON)

	deleteProduct(t, h, testCategory, productId)

	expectedCategory := makeExpectedResultJSON(t, testCategory)
	spy := spyWriter{}
	validateSavedCategory(t, h, spy, testCategory, expectedCategory)

	checkThatTheProductIsNotFoundOnTheCategory(t, h, testCategory, testProductJSON)
}

func makeExpectedResultJSON(t *testing.T, testCategory category.DTO) []byte {
	testCategory.Products = product.Products{}
	expectedCategory, err := json.Marshal(testCategory)
	if err != nil {
		t.Fatalf("could not encode test category: %s", err.Error())
	}
	return expectedCategory
}

func deleteProduct(t *testing.T, h http.Handler, testCategory category.DTO, productId product.Id) {
	deleteProductRequest, err := http.NewRequest(http.MethodDelete,
		fmt.Sprintf("/category?name=%s&product_id=%s", testCategory.Name, productId),
		http.NoBody)
	if err != nil {
		t.Fatalf("could not create request to delete category: %s", err.Error())
	}
	h.ServeHTTP(&spyWriter{}, deleteProductRequest)
}

func performRequestToGetCategoryByName(t *testing.T, h http.Handler, spy *spyWriter, categoryName category.Name) {
	request := makeRequestToGetCategoryByName(t, categoryName)
	h.ServeHTTP(spy, request)
}

func validateExpectedCategoryFromStubPersistenceInResponse(t *testing.T, spy *spyWriter, expectedCategory string) {
	expectedCategoryJSON, err := json.Marshal(expectedCategories[category.Name(expectedCategory)])
	if err != nil {
		t.Fatalf("could not parse expected category: %s", err.Error())
	}
	if fmt.Sprintf("%s", spy.written) != fmt.Sprintf("%s\n", expectedCategoryJSON) {
		t.Fatalf("did not retrieve expected category, \nexpected '%s'\nreceived '%s'\n", expectedCategoryJSON, spy.written)
	}
}

func makeRequestToGetCategoryByName(t *testing.T, categoryName category.Name) *http.Request {

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/category?name=\"%s\"", categoryName), http.NoBody)
	if err != nil {
		t.Fatalf("could not create request: %s", err.Error())
	}
	return request
}

func makeStubPersistenceHandler() http.Handler {
	testPersistence := stubPersistence{}
	h := MakeHandler(testPersistence)
	return h
}

func makePOSTRequest(t *testing.T, marshal []byte) *http.Request {
	request, err := http.NewRequest(http.MethodPost, "/category", bytes.NewReader(marshal))
	if err != nil {
		t.Fatalf("could not create request: %s", err.Error())
	}
	return request
}

func createTestCategory(t *testing.T) (category.DTO, []byte) {
	const testCategoryName = "Products for Astronauts"
	testCategory := category.DTO{
		Name: testCategoryName,
	}
	categoryJSON, err := json.Marshal(testCategory)
	if err != nil {
		t.Fatalf("could not encode (categoryJSON) request: %s", err.Error())
	}
	return testCategory, categoryJSON
}

func validateSavedCategory(t *testing.T, h http.Handler, spy spyWriter, testCategory category.DTO, expectedCategoryJSON []byte) {
	performRequestToGetCategoryByName(t, h, &spy, testCategory.Name)
	if fmt.Sprintf("%s", spy.written) != fmt.Sprintf("%s\n", expectedCategoryJSON) {
		t.Fatalf("did not retrieve expected category, \nexpected '%s'\nreceived '%s'\n", expectedCategoryJSON, spy.written)
	}
}

func postNewCategory(t *testing.T) (http.Handler, category.DTO, []byte) {
	h := MakeHandler(&fakePersistence{})
	postSpy := spyWriter{}
	testCategory, expectedCategoryJSON := createTestCategory(t)

	makeCategoryRequest := makePOSTRequest(t, expectedCategoryJSON)
	h.ServeHTTP(&postSpy, makeCategoryRequest)
	return h, testCategory, expectedCategoryJSON
}

func addProductToCategory(t *testing.T, testCategory category.DTO, h http.Handler) ([]byte, product.Id) {
	testProduct, err := product.New("test product", 1000, "a test product")
	if err != nil {
		t.Fatalf("could not create test product: %s", err.Error())
	}

	testProductJSON, err := json.Marshal(testProduct)
	addProductToCategoryRequest, err := http.NewRequest(http.MethodPut,
		fmt.Sprintf("/category?name=\"%s\"", testCategory.GetName()),
		bytes.NewReader(testProductJSON))
	if err != nil {
		t.Fatalf("could not create request to put a product in a category: %s", err.Error())
	}
	putSpy := spyWriter{}
	h.ServeHTTP(&putSpy, addProductToCategoryRequest)
	return testProductJSON, testProduct.GetId()
}

func checkThatTheProductIsFoundOnTheCategory(t *testing.T, h http.Handler, testCategory category.DTO, testProductJSON []byte) {
	newSpy := spyWriter{}
	performRequestToGetCategoryByName(t, h, &newSpy, testCategory.Name)
	if !strings.Contains(fmt.Sprintf("%s", newSpy.written), string(testProductJSON)) {
		t.Fatalf("product not found in category")
	}
}

func checkThatTheProductIsNotFoundOnTheCategory(t *testing.T, h http.Handler, testCategory category.DTO, testProductJSON []byte) {
	newSpy := spyWriter{}
	performRequestToGetCategoryByName(t, h, &newSpy, testCategory.Name)
	if strings.Contains(fmt.Sprintf("%s", newSpy.written), string(testProductJSON)) {
		t.Fatalf("product not deleted from category")
	}
}

func deleteCategory(t *testing.T, h http.Handler, name category.Name) {
	deleteRequest, err := http.NewRequest(http.MethodDelete, string("/category?name="+name), http.NoBody)
	if err != nil {
		t.Fatalf("could not create request to delete category: %s", err.Error())
	}
	h.ServeHTTP(&spyWriter{}, deleteRequest)
}

func validateCategoryDeletion(t *testing.T, h http.Handler, testCategory category.DTO, expectedCategoryJSON []byte) {
	newSpy := spyWriter{}
	performRequestToGetCategoryByName(t, h, &newSpy, testCategory.Name)
	if bytes.Contains(newSpy.written, expectedCategoryJSON) {
		t.Fatalf("category not deleted")
	}
}

type spyWriter struct {
	written []byte
	header  int
}

func (s *spyWriter) Header() http.Header {
	//TODO implement me
	panic("implement me")
}

func (s *spyWriter) Write(w []byte) (int, error) {
	s.written = append(s.written, w...)
	return len(w), nil
}

func (s *spyWriter) WriteHeader(h int) {
	s.header = h
}

type stubPersistence struct {
}

func (s stubPersistence) DeleteCategory(name category.Name) {
	//TODO implement me
	panic("implement me")
}

func (s stubPersistence) SaveCategory(dto category.DTO) {
	//TODO implement me
	panic("implement me")
}

func (s stubPersistence) GetCategories() category.Categories {
	return expectedCategories
}

type fakePersistence struct {
	categories category.Categories
}

func (f *fakePersistence) DeleteCategory(name category.Name) {
	delete(f.categories, name)
}

func (f *fakePersistence) SaveCategory(dto category.DTO) {
	if f.categories == nil {
		f.categories = make(category.Categories)
	}
	f.categories[dto.Name] = dto

}

func (f *fakePersistence) GetCategories() category.Categories {
	return f.categories
}

var expectedCategories category.Categories = category.Categories{}
var expectedCategoriesJSON []byte

func init() {
	expectedCategories["cool products"] = category.DTO{
		Name:     "cool products",
		Products: nil,
	}
	expectedCategories["hot products"] = category.DTO{
		Name:     "hot products",
		Products: nil,
	}
	expectedCategories["amazing products"] = category.DTO{
		Name:     "amazing products",
		Products: nil,
	}
	expectedCategoriesJSON, _ = json.Marshal(expectedCategories)
	expectedCategoriesJSON = append(expectedCategoriesJSON, '\n')
}
