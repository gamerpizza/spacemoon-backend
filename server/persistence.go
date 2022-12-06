package main

import (
	"errors"
	"spacemoon/login"
	"spacemoon/product"
	"spacemoon/product/category"
	"time"
)

type temporaryProductPersistence struct {
	savedProducts product.Products
}

func (t *temporaryProductPersistence) DeleteProduct(id product.Id) error {
	delete(t.savedProducts, id)
	return nil
}

func (t *temporaryProductPersistence) GetProducts() (product.Products, error) {
	return t.savedProducts, nil
}

func (t *temporaryProductPersistence) SaveProduct(p product.Product) error {
	if t.savedProducts == nil {
		t.savedProducts = make(product.Products)
	}
	t.savedProducts[p.GetId()] = p.DTO()
	return nil
}

type temporaryCategoryPersistence struct {
	categories category.Categories
}

func (t *temporaryCategoryPersistence) DeleteCategory(name category.Name) {
	delete(t.categories, name)
}

func (t *temporaryCategoryPersistence) SaveCategory(dto category.DTO) {
	if t.categories == nil {
		t.categories = make(category.Categories)
	}
	t.categories[dto.Name] = dto
}

func (t *temporaryCategoryPersistence) GetCategories() category.Categories {
	return t.categories
}

type temporaryLoginPersistence struct {
	users  map[login.User]login.Password
	tokens login.Credentials
}

func (t *temporaryLoginPersistence) ValidateCredentials(usr login.User, p login.Password) bool {
	if t.users[usr] == p {
		return true
	}
	return false
}

func (t *temporaryLoginPersistence) GetUser(token login.Token) (login.User, error) {
	tokenInfo, exists := t.tokens[token]
	if !exists {
		return "", errors.New("token not found")
	}
	if tokenInfo.Expiration.Before(time.Now()) {
		delete(t.tokens, token)
		return "", errors.New("token expired, deleted")
	}
	return tokenInfo.User, nil
}

func (t *temporaryLoginPersistence) SetUserToken(user login.User, token login.Token, tokenDuration time.Duration) {
	if t.tokens == nil {
		t.tokens = make(login.Credentials)
	}
	t.tokens[token] = login.TokenDetails{
		User:       user,
		Expiration: time.Now().Add(tokenDuration),
	}
}

var loginPersistence = &temporaryLoginPersistence{}

func init() {
	loginPersistence.users = make(map[login.User]login.Password)
	loginPersistence.users["admin"] = "sp4c3m00n!"
}
