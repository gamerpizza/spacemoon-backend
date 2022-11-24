package model

import (
	"time"
)

type Cart struct {
	UserID    UserID    `bson:"user_id" json:"userId" gorm:"column:user_id;primaryKey"`
	Products  []Product `bson:"products" json:"products"    gorm:"many2many:cart_products" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Price     float64   `bson:"price" json:"price"`
	CreatedAt time.Time `bson:"created_at" json:"createdAt" gorm:"autoCreateTime:nano"`
	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt" gorm:"autoUpdateTime:nano"`
	DeletedAt time.Time `bson:"deleted_at" json:"-"`
}

func (c Cart) Key() map[string]interface{} {
	return map[string]interface{}{
		"user_id": c.UserID,
	}
}

func (c Cart) AddProduct(p Product) {
	c.Products = append(c.Products, p)
}

func (c Cart) DeleteProduct(id string) {
	for i, v := range c.Products {
		if id == v.ID {
			c.Products = append(c.Products[:i], c.Products[i+1:]...)
			return
		}
	}
}
