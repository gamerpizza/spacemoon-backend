package main

import (
	"os"
	"spacemoon/login"
	"spacemoon/product"
	"spacemoon/product/category"
	"spacemoon/product/ratings"
	"testing"
	"time"
)

func TestGetLoginPersistence(t *testing.T) {
	var temp login.Persistence = getLoginPersistence()
	if _, ok := temp.(*temporaryLoginPersistence); !ok {
		t.Fatal("not the expected persistence")
	}

	err := os.Setenv(mongoHostKey, "host")
	if err != nil {
		t.Fatal(err.Error())
	}
	err = os.Setenv(mongoUserNameKey, "user")
	if err != nil {
		t.Fatal(err.Error())
	}
	err = os.Setenv(mongoPasswordKey, "pass")
	if err != nil {
		t.Fatal(err.Error())
	}
	var mongo login.Persistence = getLoginPersistence()
	if _, ok := mongo.(*mongoPersistence); !ok {
		t.Fatal("not the expected persistence")
	}
	os.Clearenv()
}

func TestLoginPersistence(t *testing.T) {
	persistences := map[string]login.Persistence{"temporary": &temporaryLoginPersistence{}, "mongo": &mongoPersistence{}}
	for k, per := range persistences {
		var u login.User = "user"
		var p login.Password = "pass"
		var tok login.Token = "token"
		per.SignUpUser(u, p)
		if !per.ValidateCredentials(u, p) {
			t.Fatalf("%s created credentials not working", k)
		}
		per.SetUserToken(u, tok, time.Hour)
		if usr, _ := per.GetUser(tok); usr != u {
			t.Fatalf("%s token not working", k)
		}
	}
}

func TestGetProductPersistence(t *testing.T) {
	var temp product.Persistence = getProductPersistence()
	if _, ok := temp.(*temporaryProductPersistence); !ok {
		t.Fatal("not the expected persistence")
	}
	err := os.Setenv(mongoHostKey, "host")
	if err != nil {
		t.Fatal(err.Error())
	}
	err = os.Setenv(mongoUserNameKey, "user")
	if err != nil {
		t.Fatal(err.Error())
	}
	err = os.Setenv(mongoPasswordKey, "pass")
	if err != nil {
		t.Fatal(err.Error())
	}
	var mongo product.Persistence = getProductPersistence()
	if _, ok := mongo.(*mongoPersistence); !ok {
		t.Fatal("not the expected persistence")
	}
	os.Clearenv()
}

func TestGetProductRatingsPersistence(t *testing.T) {
	var temp ratings.Persistence = getProductRatingsPersistence()
	if _, ok := temp.(*temporaryRatingsPersistence); !ok {
		t.Fatal("not the expected persistence")
	}
	err := os.Setenv(mongoHostKey, "host")
	if err != nil {
		t.Fatal(err.Error())
	}
	err = os.Setenv(mongoUserNameKey, "user")
	if err != nil {
		t.Fatal(err.Error())
	}
	err = os.Setenv(mongoPasswordKey, "pass")
	if err != nil {
		t.Fatal(err.Error())
	}
	var mongo ratings.Persistence = getProductRatingsPersistence()
	if _, ok := mongo.(*mongoPersistence); !ok {
		t.Fatal("not the expected persistence")
	}
	os.Clearenv()
}

func TestGetCategoryPersistence(t *testing.T) {
	var temp category.Persistence = getCategoryPersistence()
	if _, ok := temp.(*temporaryCategoryPersistence); !ok {
		t.Fatal("not the expected persistence")
	}
	err := os.Setenv(mongoHostKey, "host")
	if err != nil {
		t.Fatal(err.Error())
	}
	err = os.Setenv(mongoUserNameKey, "user")
	if err != nil {
		t.Fatal(err.Error())
	}
	err = os.Setenv(mongoPasswordKey, "pass")
	if err != nil {
		t.Fatal(err.Error())
	}
	var mongo category.Persistence = getCategoryPersistence()
	if _, ok := mongo.(*mongoPersistence); !ok {
		t.Fatal("not the expected persistence")
	}
	os.Clearenv()
}

var _ login.Persistence = &mongoPersistence{}
var _ product.Persistence = &mongoPersistence{}
var _ ratings.Persistence = &mongoPersistence{}
var _ category.Persistence = &mongoPersistence{}
