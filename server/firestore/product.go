package firestore

import "spacemoon/product"

func (p *fireStorePersistence) GetProducts() (product.Products, error) {
	collection := p.storage.Collection(productCollection)
	snapshots, err := collection.Documents(p.ctx).GetAll()
	if err != nil {
		return nil, err
	}

	products := make(product.Products)
	for _, snapshot := range snapshots {
		var prod product.Dto
		err := snapshot.DataTo(&prod)
		if err != nil {
			return nil, err
		}
		products[prod.Id] = prod
	}
	return products, nil
}

func (p *fireStorePersistence) SaveProduct(prod product.Product) error {
	collection := p.storage.Collection(productCollection)
	_, err := collection.Doc(string(prod.GetId())).Set(p.ctx, prod.DTO())
	if err != nil {
		return err
	}
	return nil
}

func (p *fireStorePersistence) DeleteProduct(id product.Id) error {
	collection := p.storage.Collection(productCollection)
	_, err := collection.Doc(string(id)).Delete(p.ctx)
	if err != nil {
		return err
	}
	return nil
}
