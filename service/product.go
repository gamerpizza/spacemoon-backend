package service

import (
	"context"
	"moonspace/model"
	"moonspace/repository"
	"moonspace/repository/types"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product interface {
	Create(p model.Product) error
	Get(cid, pid string) (model.Product, error)
	GetProductsLimit(cid string, start, end uint64) ([]model.Product, error)
	Update(p model.Product) error
	Delete(cid, pid string) error
}

type productImpl struct {
	categoryRepo types.ProductAssociatedRepository[model.Category]
}

func NewProductService(cli any, cfg types.Config) Product {
	return &productImpl{
		categoryRepo: repository.CreateProductRepository[model.Category](cli, cfg),
	}
}

func (s *productImpl) Create(p model.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	category := model.Category{
		ID: p.CategoryID,
	}
	p.ID = primitive.NewObjectID().Hex()
	p.CreatedAt = time.Now()

	return s.categoryRepo.AddProduct(ctx, &category, &p)
}

func (s *productImpl) Get(cid, pid string) (model.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	category := model.Category{
		ID: cid,
	}
	prod := model.Product{}

	err := s.categoryRepo.GetProduct(ctx, &category, pid, &prod)
	if err != nil {
		return model.Product{}, err
	}

	return prod, nil
}

func (s *productImpl) GetProductsLimit(cid string, start, end uint64) ([]model.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	category := model.Category{
		ID: cid,
	}

	err := s.categoryRepo.GetProductsLimit(ctx, &category, start, end)
	if err != nil {
		return nil, err
	}

	return category.Products, nil
}

func (s *productImpl) Update(p model.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	category := model.Category{
		ID: p.CategoryID,
	}

	p.UpdatedAt = time.Now()

	return s.categoryRepo.UpdateProduct(ctx, &category, &p)
}

func (s *productImpl) Delete(cid, pid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	category := model.Category{
		ID: cid,
	}
	product := model.Product{
		ID: pid,
	}

	return s.categoryRepo.DeleteProduct(ctx, &category, &product)
}
