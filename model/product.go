package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID          string    `bson:"_id,omitempty" json:"productId" gorm:"column:product_id" gorm:"primaryKey" gorm:"unique"`
	CategoryID  string    `bson:"category_id" json:"categoryId" gorm:"foreignKey:category_id"`
	Name        string    `bson:"name" json:"name"  gorm:"index"`
	Price       float64   `bson:"price" json:"price"   gorm:"index"`
	Rating      uint      `bson:"rating" json:"rating"  gorm:"index"`
	Image       string    `bson:"image" json:"image"`
	Description string    `bson:"description" json:"description"`
	CreatedBy   UserID    `bson:"createdBy" json:"createdBy" gorm:"index" gorm:"type:bytea"`
	CreatedAt   time.Time `bson:"created_at" json:"createdAt" gorm:"autoCreateTime:nano"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updatedAt" gorm:"autoUpdateTime:nano"`
	DeletedAt   time.Time `bson:"deleted_at" json:"-"`
}

func (p Product) Key() map[string]interface{} {
	id, _ := primitive.ObjectIDFromHex(p.ID)
	return map[string]interface{}{
		"_id": id,
	}
}
