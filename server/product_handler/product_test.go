package product_handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"spacemoon/product"
	"testing"
)

func TestHandler_ServeHTTP_Get(t *testing.T) {
	testHandler, fakeRequest, spy := setUpServeHTTPTest("/product")
	testHandler.ServeHTTP(&spy, fakeRequest)
	validateThatExpectedProductsAreRetrieved(t, spy)
}

func TestHandler_ServeHTTP_Get_OneProduct(t *testing.T) {
	testHandler, fakeRequest, spy := setUpServeHTTPTest(fmt.Sprintf("/product?id=%s", productId1))
	testHandler.ServeHTTP(&spy, fakeRequest)
	validateThatSpecificExpectedProductIsRetrieved(t, spy)
}

func TestHandler_ServeHTTP_Post_Then_Get(t *testing.T) {
	var fakePersistence Persistence = &fakePersistence{}
	testHandler := MakeHandler(fakePersistence)
	const productName = "Mars rocks"
	newProduct, err := product.New(productName, 1, "some description")
	productJson, err := json.Marshal(newProduct)
	if err != nil {
		return
	}
	fakePostRequest := httptest.NewRequest(http.MethodPost, "/product", bytes.NewReader(productJson))
	postSpy := spyWriter{}
	testHandler.ServeHTTP(&postSpy, fakePostRequest)
	if postSpy.header != http.StatusCreated {
		t.Fatalf("did not return the expected 201 status, got %d instead", postSpy.header)
	}
	var postResponseProduct product.Dto
	err = json.Unmarshal([]byte(postSpy.written), &postResponseProduct)
	if err != nil {
		t.Fatalf("could not unmarshal response: %s\n", err.Error())
	}

	fakeGetRequest := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/product?id=%s", postResponseProduct.GetId()), bytes.NewReader(productJson))
	getSpy := spyWriter{}
	testHandler.ServeHTTP(&getSpy, fakeGetRequest)
	if getSpy.header != http.StatusOK {
		t.Fatalf("did not receive a 200 status, instead: %d\n", getSpy.header)
	}

	var getResponseProduct product.Dto
	err = json.Unmarshal([]byte(getSpy.written), &getResponseProduct)
	if err != nil {
		t.Fatalf("could not unmarshal response: %s\n", err.Error())
	}

	if postResponseProduct.DTO() != getResponseProduct.DTO() {
		t.Fatalf("retrieved DTO (%+v)is not expected DTO (%+v)", getResponseProduct.DTO(), postResponseProduct.DTO())
	}
}

func validateThatExpectedProductsAreRetrieved(t *testing.T, spy spyWriter) {
	var products product.Products = make(product.Products)
	err := json.Unmarshal([]byte(spy.written), &products)
	if err != nil {
		t.Fatalf("could not unmarshal response: %s", err.Error())
	}
	if !reflect.DeepEqual(products, expectedProducts) {
		t.Fatalf("retrieved products '%+v'\n"+
			"do not match the expected products '%+v'\n", products, expectedProducts)
	}
	if spy.header != http.StatusOK {
		t.Fatalf("did not get a status 200 on GET")
	}
}

func validateThatSpecificExpectedProductIsRetrieved(t *testing.T, spy spyWriter) {
	var p product.Dto
	err := json.Unmarshal([]byte(spy.written), &p)
	if err != nil {
		t.Fatalf("could not unmarshal response: %s", err.Error())
	}
	if !reflect.DeepEqual(p, expectedProducts[productId1]) {
		t.Fatalf("retrieved product '%+v'\n"+
			"does not match the expected product '%+v'\n", p, expectedProducts[productId1])
	}
	if spy.header != http.StatusOK {
		t.Fatalf("did not get a status 200 on GET")
	}
}

func setUpServeHTTPTest(target string) (http.Handler, *http.Request, spyWriter) {
	var fakePersistence Persistence = stubPersistence{}
	testHandler := MakeHandler(fakePersistence)
	fakeRequest := httptest.NewRequest(http.MethodGet, target, http.NoBody)
	spy := spyWriter{}
	return testHandler, fakeRequest, spy
}

type spyWriter struct {
	written string
	header  int
}

func (s *spyWriter) Header() http.Header {
	//TODO implement me
	panic("implement me")
}

func (s *spyWriter) Write(bytes []byte) (int, error) {
	s.written = s.written + fmt.Sprintf("%s", bytes)
	return len(bytes), nil
}

func (s *spyWriter) WriteHeader(h int) {
	s.header = h
}

type stubPersistence struct {
}

func (s stubPersistence) GetProducts() (product.Products, error) {
	return expectedProducts, nil
}

func (s stubPersistence) SaveProduct(p product.Product) error {
	//TODO implement me
	panic("implement me")
}

func (s stubPersistence) DeleteProduct(id product.Id) error {
	//TODO implement me
	panic("implement me")
}

type fakePersistence struct {
	savedProducts product.Products
}

func (f *fakePersistence) GetProducts() (product.Products, error) {
	return f.savedProducts, nil
}

func (f *fakePersistence) SaveProduct(p product.Product) error {
	if f.savedProducts == nil {
		f.savedProducts = make(product.Products)
	}
	f.savedProducts[p.GetId()] = p.DTO()
	return nil
}

func (f *fakePersistence) DeleteProduct(id product.Id) error {
	//TODO implement me
	panic("implement me")
}

var expectedProducts = make(product.Products)

func init() {
	expectedProducts[productId1] = product.Dto{
		Name:        "product1",
		Price:       1,
		Description: "",
	}
	expectedProducts["product2-id"] = product.Dto{
		Name:        "product2",
		Price:       10,
		Description: "",
	}
}

const productId1 = "product1-id"
