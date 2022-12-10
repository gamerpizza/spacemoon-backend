package main

import (
	"errors"
	"log"
	"os"
	"spacemoon/login"
	"spacemoon/product"
	"spacemoon/product/category"
	"spacemoon/product/ratings"
	"strings"
	"time"
)

func getLoginPersistence() login.Persistence {
	mongoHost, mongoUsr, mongoPass := getMongoEnvironmentVariables()
	if checkIfMongoParametersAreNotValid(mongoHost, mongoUsr, mongoPass, "login") {
		per := &temporaryLoginPersistence{}
		//hard coded credentials
		per.users = make(map[login.User]login.Password)
		per.users["admin"] = "sp4c3m00n!"
		return per
	}
	return &mongoPersistence{}
}

func getProductPersistence() product.Persistence {
	mongoHost, mongoUsr, mongoPass := getMongoEnvironmentVariables()
	if checkIfMongoParametersAreNotValid(mongoHost, mongoUsr, mongoPass, "product") {
		return &temporaryProductPersistence{}
	}
	return &mongoPersistence{}
}

func getProductRatingsPersistence() ratings.Persistence {
	mongoHost, mongoUsr, mongoPass := getMongoEnvironmentVariables()
	if checkIfMongoParametersAreNotValid(mongoHost, mongoUsr, mongoPass, "ratings") {
		return &temporaryRatingsPersistence{}
	}
	return &mongoPersistence{}
}

func getCategoryPersistence() category.Persistence {
	mongoHost, mongoUsr, mongoPass := getMongoEnvironmentVariables()
	if checkIfMongoParametersAreNotValid(mongoHost, mongoUsr, mongoPass, "category") {
		return &temporaryCategoryPersistence{}
	}
	return &mongoPersistence{}
}

type mongoPersistence struct {
}

func (m *mongoPersistence) SignUpUser(_ login.User, _ login.Password) {
}

func (m *mongoPersistence) SetUserToken(_ login.User, _ login.Token, _ time.Duration) {

}

func (m *mongoPersistence) GetUser(_ login.Token) (login.User, error) {
	return "", nil
}

func (m *mongoPersistence) ValidateCredentials(_ login.User, _ login.Password) bool {
	return false
}

func (m *mongoPersistence) GetCategories() category.Categories {
	//TODO implement me
	panic("implement me")
}

func (m *mongoPersistence) SaveCategory(dto category.DTO) {
	//TODO implement me
	panic("implement me")
}

func (m *mongoPersistence) DeleteCategory(name category.Name) {
	//TODO implement me
	panic("implement me")
}

func (m *mongoPersistence) ReadRating(id product.Id) ratings.Rating {
	//TODO implement me
	panic("implement me")
}

func (m *mongoPersistence) SaveRating(id product.Id, rating ratings.Rating) {
	//TODO implement me
	panic("implement me")
}

func (m *mongoPersistence) GetProducts() (product.Products, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mongoPersistence) SaveProduct(p product.Product) error {
	//TODO implement me
	panic("implement me")
}

func (m *mongoPersistence) DeleteProduct(id product.Id) error {
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
	users  login.Credentials
	tokens login.Tokens
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
		t.tokens = make(login.Tokens)
	}
	t.tokens[token] = login.TokenDetails{
		User:       user,
		Expiration: time.Now().Add(tokenDuration),
	}
}

func (t *temporaryLoginPersistence) SignUpUser(u login.User, p login.Password) {
	if t.users == nil {
		t.users = make(login.Credentials)
	}
	t.users[u] = p
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
