package product

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
	var fakePersistence Persistence = stubPersistence{}
	testHandler := MakeHandler(fakePersistence)
	fakeRequest := httptest.NewRequest(http.MethodGet, "/product", http.NoBody)
	spy := spyWriter{}
	testHandler.ServeHTTP(&spy, fakeRequest)
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

func TestHandler_ServeHTTP_Get_OneProduct(t *testing.T) {
	var fakePersistence Persistence = stubPersistence{}
	testHandler := MakeHandler(fakePersistence)
	fakeRequest := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/product?id=%s", productId1), http.NoBody)
	spy := spyWriter{}
	testHandler.ServeHTTP(&spy, fakeRequest)
	var p product.Product
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

func TestHandler_ServeHTTP_Post(t *testing.T) {
	var fakePersistence Persistence = stubPersistence{}
	testHandler := MakeHandler(fakePersistence)
	newProduct, err := product.New("test-product", 1, "some description")
	marshal, err := json.Marshal(newProduct)
	if err != nil {
		return
	}
	fakeRequest := httptest.NewRequest(http.MethodPost, "/product", bytes.NewReader(marshal))
	spy := spyWriter{}
	testHandler.ServeHTTP(&spy, fakeRequest)

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

func (s stubPersistence) GetProducts() product.Products {
	return expectedProducts
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
