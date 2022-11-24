package dto

import (
	"mime/multipart"
	"moonspace/model"
)

type CategoryDto struct {
	Name  string                `form:"name"  binding:"required"`
	Image *multipart.FileHeader `form:"image" binding:"required"`
}

type ProductDto struct {
	Name        string                `form:"name" json:"name" binding:"required"`
	Price       float64               `form:"price" json:"price"  binding:"required"`
	Image       *multipart.FileHeader `form:"image" json:"image" binding:"required"`
	Description string                `form:"description" json:"description" binding:"required"`
}

func ProductDtoToModel(p ProductDto, imgPath, cid string, uid model.UserID) model.Product {
	prod := model.Product{}
	prod.Description = p.Description
	prod.Image = imgPath
	prod.Price = p.Price
	prod.Name = p.Name
	prod.CategoryID = cid
	prod.CreatedBy = uid

	return prod
}

func CategoryDtoToModel(c CategoryDto, imgPath string, uid model.UserID) model.Category {
	category := model.Category{}
	category.CreatedBy = uid
	category.Name = c.Name
	category.Image = imgPath

	return category
}
