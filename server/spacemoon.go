package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"spacemoon/login"
	"spacemoon/product"
	"spacemoon/product/category"
	"spacemoon/server/category_handler"
	"spacemoon/server/product_handler"
	"time"
)

func main() {
	setupHandlers()
	listenAndServe()
}

func setupHandlers() {
	log.Default().Print("starting spacemoon server 🚀")
	log.Default().Print("registering server handlers...")
	http.Handle("/product", product_handler.MakeHandler(&temporaryProductPersistence{}))
	http.Handle("/category", category_handler.MakeHandler(&temporaryCategoryPersistence{}))
	http.Handle("/login", login.NewHandler(loginPersistence, time.Hour))
	log.Default().Print("handler registration done, ready for takeoff")
}

func listenAndServe() {
	log.Default().Print(getRandomSpaceQuote())
	log.Default().Printf("listening on port %d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Default().Fatalf("error while performing listen and serve: %s", err.Error())
	}
}

func getRandomSpaceQuote() string {
	rand.Seed(time.Now().Unix())
	quote := rand.Intn(2)
	switch quote {
	case 0:
		return "“The stars don't look bigger, but they do look brighter.” ― Sally Ride"
	case 1:
		return "“I see Earth! It is so beautiful.” ― Yuri Gagarin"
	default:
		return ""
	}
}

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

const port = 1234

var loginPersistence = &temporaryLoginPersistence{}

func init() {
	loginPersistence.users = make(map[login.User]login.Password)
	loginPersistence.users["admin"] = "sp4c3m00n!"
}
