package category

import (
	"errors"
	"spacemoon/product"
	"strings"
)

// New instantiates a Category. It will return an EmptyNameError if the Name provided is empty.
func New(n Name) (Category, error) {
	if strings.TrimSpace(string(n)) == "" {
		return nil, EmptyNameError
	}
	return &DTO{Name: n}, nil
}

// Category is defined as a group of product.Product. A product.Product can be on more than one category
// at the same time
type Category interface {
	GetName() Name
	GetProducts() product.Products
	AddProduct(product.Product) product.Id
	DTO() DTO
}

type DTO struct {
	Name     Name             `json:"name"`
	Products product.Products `json:"products"`
}

func (d *DTO) DTO() DTO {
	return *d
}

func (d *DTO) AddProduct(p product.Product) product.Id {
	if d.Products == nil {
		d.Products = make(product.Products)
	}
	d.Products[p.GetId()] = p.DTO()
	return p.GetId()
}

func (d *DTO) GetProducts() product.Products {
	return d.Products
}

// GetName returns the name of the category
func (d *DTO) GetName() Name {
	return d.Name
}

func (d *DTO) DeleteProduct(id product.Id) {
	delete(d.Products, id)
}

type Name string
type Categories map[Name]DTO

var EmptyNameError = errors.New("category name cannot be empty")
