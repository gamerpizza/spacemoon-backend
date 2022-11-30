package service

import (
	"context"
	"moonspace/model"
	"moonspace/repository"
	"moonspace/repository/types"
	"time"
)

type Product interface {
	Create(p model.Product) error
	Get(pid string) (model.Product, error)
	GetProductsLimit(cid string, start, end uint64) ([]model.Product, error)
	Update(pid string, p model.Product) error
	Delete(cid, pid string) error
}

type productImpl struct {
	productRepo types.Repository[model.Product]
}

func NewProductService(cli any, cfg types.Config) Product {
	return &productImpl{
		productRepo: repository.CreateRepository[model.Product](cli, cfg),
	}
}

func (s *productImpl) Create(p model.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	p.CreatedAt = time.Now()
	return s.productRepo.Add(ctx, p)
}

func (s *productImpl) Get(pid string) (model.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	p := model.Product{ID: pid}

	err := s.productRepo.Get(ctx, p.Key(), &p)
	if err != nil {
		return model.Product{}, err
	}

	return p, nil
}

func (s *productImpl) GetProductsLimit(cid string, start, end uint64) ([]model.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	p := make([]model.Product, 0)

	err := s.productRepo.GetProductLimit(ctx, cid, start, end, &p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *productImpl) Update(pid string, p model.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	p.UpdatedAt = time.Now()

	return s.productRepo.UpdateProduct(ctx, p.CategoryID, pid, p)
}

func (s *productImpl) Delete(cid, pid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	product := model.Product{
		ID: pid,
	}

	return s.productRepo.DeleteProduct(ctx, cid, product)
}
