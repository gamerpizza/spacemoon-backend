// Package product manages the core business logic of what a basic product is and how it is defined.
package product

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"spacemoon/login"
	"strings"
)

// New instantiates a new product with the main required attributes set (Name, Price and Description). Because
// the rating of the product is a complex responsibility by itself, it makes sense to have a Rater somewhere else,
// and leave the rating out of the product itself.
// If the provided Name for the new Product is empty, New will return nil and an error.
func New(n Name, p Price, d Description, seller login.UserName) (Product, error) {
	if strings.TrimSpace(string(n)) == "" {
		return nil, EmptyNameError
	}
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("could not generate uuid to create new product: %w", err)
	}
	return &Dto{
		Name:        n,
		Price:       p,
		Description: d,
		Id:          Id(id.String()),
		Seller:      seller,
	}, nil
}

// Product is defined as an interface for encapsulation purposes
type Product interface {
	GetName() Name
	GetId() Id
	GetPrice() Price
	GetDescription() Description
	DTO() Dto
	GetType() Type
	SetType(Type)
}

// Dto is a representation of a product, with a given Price, Name, Rating and Description, used for data transfer.
type Dto struct {
	Name        Name           `json:"name"`
	Price       Price          `json:"price"`
	Description Description    `json:"description"`
	Id          Id             `json:"id"`
	Type        Type           `json:"type"`
	Seller      login.UserName `json:"seller"`
}

func (d *Dto) SetType(t Type) {
	d.Type = t
}

func (d *Dto) GetType() Type {
	return d.Type
}

func (d *Dto) GetId() Id {
	return d.Id
}

// DTO returns the DTO form of a product, structurally similar to what a String() method would do (but it returns a DTO
// and not a string).
func (d *Dto) DTO() Dto {
	return *d
}

func (d *Dto) GetName() Name {
	return d.Name
}

func (d *Dto) GetPrice() Price {
	return d.Price
}

func (d *Dto) GetDescription() Description {
	return d.Description
}

// Name is the name of the product. I may be non-unique, as the product uses a UUID to identify different SKUs
// and different instances of the same product with different characteristics.
type Name string

// Id returns is Product unique identifier (UUID). We are not using the Name as Id because two products could have
// the same Name with different characteristics, and they would still need to be treated as two different SKUs (two
// unique products).
type Id string

// Price is an int to avoid floating comma errors when working with fractional amounts. It makes more sense to use an
// integer number of the smallest monetary unit possible (cents, satoshis, etc.) than to use a floating point.
type Price int

// Description is an open field to add any information about the product that won't participate in the business logic,
// as a form of "notes" or "extra information" normally expected to be displayed somewhere in the front-end.
type Description string

// Type is a short description/name for the type of product
type Type string

// Products represents the product.Product contents of a category
type Products map[Id]Dto

// EmptyNameError is returned when trying to create a product with an empty name.
var EmptyNameError = errors.New("name cannot be empty")

// Persistence is used, as expected, to write and read, to be able to save information.
type Persistence interface {
	GetProducts() (Products, error)
	SaveProduct(Product) error
	DeleteProduct(Id) error
}
