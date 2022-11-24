package model

type Entity interface {
	Key() map[string]interface{}
}

type ProductEntity interface {
	AddProduct(p Product)
	DeleteProduct(id string)
	Entity
}
