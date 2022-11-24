package model

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	//Unique identificator
	ID string `bson:"_id,omitempty" json:"categoryId" gorm:"column:category_id" gorm:"primaryKey" gorm:"unique"`
	//Category name
	Name string `bson:"name" json:"name" gorm:"index"`
	//Category image
	Image string `bson:"image" json:"image"`
	//Category products
	Products []Product `bson:"products" json:"products"    gorm:"foreignKey:ID" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	//User that has created the category
	CreatedBy UserID `bson:"createdBy" json:"createdBy" gorm:"index" gorm:"type:bytea"`
	//Time of creation
	CreatedAt time.Time `bson:"created_at" json:"createdAt" gorm:"autoCreateTime:nano"`
	//Time of last update
	UpdatedAt time.Time `bson:"updated_at" json:"updatedAt" gorm:"autoUpdateTime:nano"`
	//Time of logical deletion
	DeletedAt time.Time `bson:"deleted_at" json:"-"`
}

func (cat Category) Key() map[string]interface{} {
	id, _ := primitive.ObjectIDFromHex(cat.ID)
	return map[string]interface{}{
		"_id": id,
	}
}

func (cat Category) AddProduct(p Product) {
	cat.Products = append(cat.Products, p)
}

func (cat Category) DeleteProduct(id string) {
	for i, v := range cat.Products {
		if id == v.ID {
			cat.Products = append(cat.Products[:i], cat.Products[i+1:]...)
			return
		}
	}
}
