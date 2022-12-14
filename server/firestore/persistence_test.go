package firestore

import (
	"context"
	"spacemoon/product/category"
	"spacemoon/product/ratings"
	"testing"
)

func TestGetPersistence(t *testing.T) {
	ctx := context.TODO()

	var rp ratings.Persistence
	rp, _ = GetPersistence(ctx)
	rp.ReadRating("")

	var cp category.Persistence
	cp, _ = GetPersistence(ctx)
	cp.GetCategories()
}
