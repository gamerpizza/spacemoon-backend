package main

import (
	"context"
	"errors"
	"log"
	"os"
	"spacemoon/login"
	"spacemoon/product"
	"spacemoon/product/category"
	"spacemoon/product/ratings"
	"spacemoon/server/firestore"
	"strings"
	"time"
)

func getLoginPersistence() login.Persistence {
	mongoHost, mongoUsr, mongoPass := getMongoEnvironmentVariables()
	if checkIfMongoParametersAreNotValid(mongoHost, mongoUsr, mongoPass, "login") {
		per := &temporaryLoginPersistence{}
		//hard coded credentials
		per.users = make(map[login.UserName]login.Password)
		per.users["admin"] = "sp4c3m00n!"
		return per
	}
	per, _ := firestore.GetPersistence(context.Background())
	return per
}

func getProductPersistence() product.Persistence {
	mongoHost, mongoUsr, mongoPass := getMongoEnvironmentVariables()
	if checkIfMongoParametersAreNotValid(mongoHost, mongoUsr, mongoPass, "product") {
		return &temporaryProductPersistence{}
	}
	return &googleCloudPersistence{}
}

func getProductRatingsPersistence() ratings.Persistence {
	mongoHost, mongoUsr, mongoPass := getMongoEnvironmentVariables()
	if checkIfMongoParametersAreNotValid(mongoHost, mongoUsr, mongoPass, "ratings") {
		return &temporaryRatingsPersistence{}
	}
	return &googleCloudPersistence{}
}

func getCategoryPersistence() category.Persistence {
	mongoHost, mongoUsr, mongoPass := getMongoEnvironmentVariables()
	if checkIfMongoParametersAreNotValid(mongoHost, mongoUsr, mongoPass, "category") {
		return &temporaryCategoryPersistence{}
	}
	return &googleCloudPersistence{}
}

type googleCloudPersistence struct {
}

func (m *googleCloudPersistence) SignUpUser(_ login.UserName, _ login.Password) error {
	return nil
}

func (m *googleCloudPersistence) DeleteUser(name login.UserName) error {
	//TODO implement me
	panic("implement me")
}

func (m *googleCloudPersistence) SetUserToken(user login.UserName, token login.Token, expirationTime time.Duration) error {

	return nil
}

func (m *googleCloudPersistence) GetUser(_ login.Token) (login.UserName, error) {
	return "", nil
}

func (m *googleCloudPersistence) ValidateCredentials(_ login.UserName, _ login.Password) bool {
	return false
}

func (m *googleCloudPersistence) GetCategories() category.Categories {
	//TODO implement me
	panic("implement me")
}

func (m *googleCloudPersistence) SaveCategory(dto category.DTO) {
	//TODO implement me
	panic("implement me")
}

func (m *googleCloudPersistence) DeleteCategory(name category.Name) {
	//TODO implement me
	panic("implement me")
}

func (m *googleCloudPersistence) ReadRating(id product.Id) ratings.Rating {
	//TODO implement me
	panic("implement me")
}

func (m *googleCloudPersistence) SaveRating(id product.Id, rating ratings.Rating) {
	//TODO implement me
	panic("implement me")
}

func (m *googleCloudPersistence) GetProducts() (product.Products, error) {
	//TODO implement me
	panic("implement me")
}

func (m *googleCloudPersistence) SaveProduct(p product.Product) error {
	//TODO implement me
	panic("implement me")
}

func (m *googleCloudPersistence) DeleteProduct(id product.Id) error {
	//TODO implement me
	panic("implement me")
}

func checkIfMongoParametersAreNotValid(mongoHost, mongoUsr, mongoPass, persistenceType string) bool {
	isInvalid := strings.TrimSpace(mongoHost) == "" || strings.TrimSpace(mongoUsr) == "" || strings.TrimSpace(mongoPass) == ""
	if isInvalid {
		log.Default().Printf("using temporary persistence for %s", persistenceType)
		return isInvalid
	}
	log.Default().Printf("using mongo persistence for %s", persistenceType)
	return isInvalid
}

func getMongoEnvironmentVariables() (string, string, string) {
	mongoHost := os.Getenv(mongoHostKey)
	mongoUsr := os.Getenv(mongoUserNameKey)
	mongoPass := os.Getenv(mongoPasswordKey)
	return mongoHost, mongoUsr, mongoPass
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
	users  Credentials
	tokens login.Tokens
}

func (t *temporaryLoginPersistence) DeleteUser(name login.UserName) error {
	//TODO implement me
	panic("implement me")
}

func (t *temporaryLoginPersistence) ValidateCredentials(usr login.UserName, p login.Password) bool {
	if t.users[usr] == p {
		return true
	}
	return false
}

func (t *temporaryLoginPersistence) GetUser(token login.Token) (login.UserName, error) {
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

func (t *temporaryLoginPersistence) SetUserToken(user login.UserName, token login.Token, expirationTime time.Duration) error {
	if t.tokens == nil {
		t.tokens = make(login.Tokens)
	}
	t.tokens[token] = login.TokenDetails{
		User:       user,
		Expiration: time.Now().Add(tokenDuration),
	}
	return nil
}

func (t *temporaryLoginPersistence) SignUpUser(u login.UserName, p login.Password) error {
	if t.users == nil {
		t.users = make(Credentials)
	}
	t.users[u] = p
	return nil
}

type temporaryRatingsPersistence struct {
	r ratings.Ratings
}

func (t *temporaryRatingsPersistence) ReadRating(id product.Id) ratings.Rating {
	return t.r[id]
}

func (t *temporaryRatingsPersistence) SaveRating(id product.Id, rating ratings.Rating) {
	if t.r == nil {
		t.r = make(ratings.Ratings)
	}
	t.r[id] = rating
}

const mongoHostKey = "MONGO_HOST"
const mongoUserNameKey = "MONGO_USER"
const mongoPasswordKey = "MONGO_PASS"

type Credentials map[login.UserName]login.Password
