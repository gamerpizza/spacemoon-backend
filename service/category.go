package service

import (
	"context"
	"moonspace/model"
	"moonspace/repository"
	"moonspace/repository/types"
	"time"
)

type Category interface {
	Create(cat model.Category) error
	Delete(cid string) error
	Get(cid string) (model.Category, error)
	GetLimit(start, end uint64) ([]model.Category, error)
	Update(cat model.Category) error
}

type categoryImpl struct {
	categoryRepo types.Repository[model.Category]
}

func NewCategoryService(cli any, cfg types.Config) Category {
	return &categoryImpl{
		categoryRepo: repository.CreateRepository[model.Category](cli, cfg),
	}
}

func (c *categoryImpl) Create(cat model.Category) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cat.CreatedAt = time.Now()
	cat.Products = make([]model.Product, 0)

	return c.categoryRepo.Add(ctx, cat)
}

func (c *categoryImpl) Get(cid string) (model.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cat := model.Category{ID: cid}

	err := c.categoryRepo.Get(ctx, cat.Key(), &cat)
	if err != nil {
		return model.Category{}, err
	}

	return cat, nil
}

func (c *categoryImpl) GetLimit(start, end uint64) ([]model.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cats := make([]model.Category, 0)

	err := c.categoryRepo.GetLimit(ctx, start, end, &cats)
	if err != nil {
		return nil, err
	}

	return cats, nil
}

func (c *categoryImpl) Update(cat model.Category) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cat.UpdatedAt = time.Now()

	return c.categoryRepo.Update(ctx, cat)
}

func (c *categoryImpl) Delete(cid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cat := model.Category{
		ID: cid,
	}

	return c.categoryRepo.Delete(ctx, cat)
}
