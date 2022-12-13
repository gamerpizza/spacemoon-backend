package firestore

import (
	"context"
	"os"
	"spacemoon/login"
	"spacemoon/product"
	"spacemoon/product/category"
	"spacemoon/product/ratings"
	"testing"
)

func TestGetPersistence(t *testing.T) {
	ctx := context.TODO()

	err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/home/hari/Dev/spacemoon/spacemoon-backend/server/firestore/global-pagoda-368419-1b8a2cf2d395.json")
	if err != nil {
		t.Fatal(err.Error())
	}
	lp, err := GetPersistence(ctx)
	if err != nil {
		t.Fatal(err.Error())
	}
	lp.ValidateCredentials("", "")

	var pp product.Persistence
	pp, _ = GetPersistence(ctx)
	_, _ = pp.GetProducts()

	var rp ratings.Persistence
	rp, _ = GetPersistence(ctx)
	rp.ReadRating("")

	var cp category.Persistence
	cp, _ = GetPersistence(ctx)
	cp.GetCategories()
}

func TestLoginPersistence(t *testing.T) {
	ctx := context.TODO()
	var lp login.Persistence
	lp, err := GetPersistence(ctx)
	defer func(persistence *fireStorePersistence) {
		err := persistence.Close()
		if err != nil {
			t.Fatal(err.Error())
		}
	}(lp.(*fireStorePersistence))
	if err != nil {
		t.Fatal(err.Error())
	}

	const testUser = "test-user"
	const testPassword = "test-pass"
	err = lp.SignUpUser(testUser, testPassword)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !lp.ValidateCredentials(testUser, testPassword) {
		t.Fatal("expected user not found")
	}

	err = lp.DeleteUser(testUser)
	if err != nil {
		t.Fatal(err.Error())
	}
	if lp.ValidateCredentials(testUser, testPassword) {
		t.Fatal("unexpected user found")
	}
}
