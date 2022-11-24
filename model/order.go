package model

import (
	"time"
)

type OrderStatus uint8

const (
	OrderStatusPending   OrderStatus = 0
	OrderStatusApprooved OrderStatus = 1
	OrderStatusCanceled  OrderStatus = 2
)

type Order struct {
	OrderID   string      `bson:"order_id" json:"orderId" gorm:"column:order_id" gorm:"primaryKey" gorm:"unique"`
	Cart      Cart        `bson:"cart" json:"cart" gorm:"foreignKey:UserID"`
	CreatedBy UserID      `bson:"createdBy" json:"createdBy" gorm:"index"`
	Status    OrderStatus `bson:"status" json:"status"`
	CreatedAt time.Time   `bson:"created_at" json:"createdAt" gorm:"autoCreateTime:nano"`
	UpdatedAt time.Time   `bson:"updated_at" json:"updatedAt" gorm:"autoUpdateTime:nano"`
	DeletedAt time.Time   `bson:"deleted_at" json:"-"`
}

func (o Order) Key() map[string]interface{} {
	return map[string]interface{}{
		"order_id": o.OrderID,
	}
}
