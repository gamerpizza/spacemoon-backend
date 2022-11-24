package test_helpers

import (
	"moonspace/model"
	"time"

	"github.com/google/uuid"
)

func CreateTestData() ([]model.Product, model.Cart, model.Category, model.Order, model.UserID) {
	uid := uuid.New()
	timeNow := time.Now()

	category := model.Category{
		ID:        "1",
		Name:      "Q",
		Image:     "",
		CreatedBy: model.UUIDToUserID(uid),
		CreatedAt: timeNow,
		Products:  make([]model.Product, 0),
	}

	products := []model.Product{
		{
			ID:         "1",
			Name:       "X",
			CategoryID: "1",
			Price:      100,
			Rating:     5,
			Image:      "",
			CreatedBy:  model.UUIDToUserID(uid),
			CreatedAt:  timeNow,
		},
		{
			ID:         "2",
			Name:       "Z",
			CategoryID: "1",
			Price:      200,
			Rating:     5,
			Image:      "",
			CreatedBy:  model.UUIDToUserID(uid),
			CreatedAt:  timeNow,
		},
		{
			ID:         "3",
			Name:       "Y",
			CategoryID: "1",
			Price:      300,
			Rating:     5,
			Image:      "",
			CreatedBy:  model.UUIDToUserID(uid),
			CreatedAt:  timeNow,
		},
	}

	cart := model.Cart{
		UserID:    model.UUIDToUserID(uid),
		Products:  make([]model.Product, 0),
		Price:     0,
		CreatedAt: timeNow,
	}

	order := model.Order{
		OrderID:   "1",
		Cart:      cart,
		CreatedBy: model.UUIDToUserID(uid),
		CreatedAt: timeNow,
	}

	return products, cart, category, order, model.UUIDToUserID(uid)
}
