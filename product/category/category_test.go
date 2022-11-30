package category

import (
	"errors"
	"github.com/google/uuid"
	"reflect"
	"spacemoon/product"
	"testing"
)

func TestCategory_Name(t *testing.T) {
	var expectedName Name = "test-name"
	var testCategory Category
	var err error
	testCategory, err = New(expectedName)
	if wasThrown(err) {
		t.Fatalf("error instantiatin new category: %s", err.Error())
	}
	var retrievedName Name = testCategory.GetName()
	if retrievedName != expectedName {
		t.Fatalf("invalid category name retreived, \n"+
			"expected '%s' \n"+
			"retrieved '%s'\n", expectedName, retrievedName)
	}
}

func TestCategory_DTO(t *testing.T) {
	var expectedName Name = "test-name"
	var testCategory Category
	var err error
	testCategory, err = New(expectedName)
	if wasThrown(err) {
		t.Fatalf("error instantiatin new category: %s", err.Error())
	}
	var retrievedDTO DTO = testCategory.DTO()
	expectedDto := DTO{
		Name:     expectedName,
		Products: nil,
	}
	if areDifferent(retrievedDTO, expectedDto) {
		t.Fatalf("invalid DTO retrieved, \n "+
			"expected '%+v'\n"+
			"retrieved '%+v'\n", expectedDto, retrievedDTO)
	}
}

func TestCategory_NameCannotBeEmpty(t *testing.T) {
	var emptyName Name = ""
	var err error
	_, err = New(emptyName)
	if wasNotThrown(err) {
		t.Fatalf("no error thrown on empty category name")
	}
	if !errors.Is(err, EmptyNameError) {
		t.Fatalf("invalid error thrown on empty category name, \n"+
			"expected '%s'\n"+
			"received '%s\n", EmptyNameError.Error(), err.Error())
	}
}

func TestCategory_Products(t *testing.T) {
	testCategory, err := New("test-category")
	if wasThrown(err) {
		t.Fatalf("could not create new category: %s", err.Error())
	}

	var fakeProduct product.Product = stubProduct{}
	testCategory.AddProduct(fakeProduct)
	var retrievedProducts = testCategory.GetProducts()
	if _, exists := retrievedProducts[expectedProductID]; !exists {
		t.Fatalf("expected product not found in category")
	}
}

func areDifferent(retrievedDTO DTO, expectedDto DTO) bool {
	return retrievedDTO.Name != expectedDto.Name || !reflect.DeepEqual(retrievedDTO.Products, expectedDto.Products)
}

func wasNotThrown(err error) bool {
	return err == nil
}

func wasThrown(err error) bool {
	return err != nil
}

type stubProduct struct {
}

func (f stubProduct) GetId() product.Id {
	return expectedProductID
}

func (f stubProduct) GetName() product.Name {
	return expectedProductName
}

func (f stubProduct) GetPrice() product.Price {
	//TODO implement me
	panic("implement me")
}

func (f stubProduct) GetDescription() product.Description {
	//TODO implement me
	panic("implement me")
}

func (f stubProduct) DTO() product.Dto {
	return product.Dto{}
}

const expectedProductName = "some product"

var expectedProductID = product.Id(uuid.New().String())
