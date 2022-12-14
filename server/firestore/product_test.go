package firestore

import (
	"context"
	"spacemoon/product"
	"testing"
)

func TestProductPersistence(t *testing.T) {
	ctx := context.TODO()

	var pp product.Persistence
	pp, err := GetPersistence(ctx)
	defer func(persistence *fireStorePersistence) {
		_ = persistence.Close()
	}(pp.(*fireStorePersistence))
	if err != nil {
		t.Fatal(err.Error())
	}
	const productName = "product-name"
	p, err := product.New(productName, 100, "")
	if err != nil {
		t.Fatal(err.Error())
	}
	id := p.GetId()
	err = pp.SaveProduct(p)
	if err != nil {
		t.Fatal(err.Error())
	}

	products, err := pp.GetProducts()
	if err != nil {
		t.Fatal(err.Error())
	}
	if _, exists := products[id]; !exists {
		t.Fatal("expected product not found")
	}

	err = pp.DeleteProduct(id)
	if err != nil {
		t.Fatal(err.Error())
	}
}
