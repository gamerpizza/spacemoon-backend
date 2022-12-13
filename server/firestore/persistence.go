package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"spacemoon/login"
	"spacemoon/product"
	"spacemoon/product/category"
	"spacemoon/product/ratings"
	"time"
)

type Persistence interface {
	login.Persistence
	product.Persistence
	ratings.Persistence
	category.Persistence
	Close() error
}

func GetPersistence(ctx context.Context) (Persistence, error) {
	const projectID = "global-pagoda-368419"
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("could not create firestore client: %w", err)
	}

	return &fireStorePersistence{ctx: ctx, storage: client}, nil
}

type fireStorePersistence struct {
	storage *firestore.Client
	ctx     context.Context
}

func (p *fireStorePersistence) Close() error {
	err := p.storage.Close()
	if err != nil {
		return err
	}
	return nil
}

func (p *fireStorePersistence) GetCategories() category.Categories {
	return nil
}

func (p *fireStorePersistence) SaveCategory(dto category.DTO) {
	//TODO implement me
	panic("implement me")
}

func (p *fireStorePersistence) DeleteCategory(name category.Name) {
	//TODO implement me
	panic("implement me")
}

func (p *fireStorePersistence) ReadRating(_ product.Id) ratings.Rating {
	return ratings.Rating{}
}

func (p *fireStorePersistence) SaveRating(id product.Id, rating ratings.Rating) {
	//TODO implement me
	panic("implement me")
}

func (p *fireStorePersistence) SetUserToken(user login.UserName, token login.Token, duration time.Duration) error {
	collection := p.storage.Collection(loginTokensCollection)
	_, err := collection.Doc(string(token)).Set(p.ctx, login.Credential{
		Token: token,
		TokenDetails: login.TokenDetails{
			User:       user,
			Expiration: time.Now().Add(duration),
		},
	})
	if err != nil {
		return fmt.Errorf("could not ser user token: %w", err)
	}
	return nil
}

func (p *fireStorePersistence) GetUser(token login.Token) (login.UserName, error) {
	collection := p.storage.Collection(loginTokensCollection)
	get, err := collection.Doc(string(token)).Get(p.ctx)
	if err != nil {
		return "", fmt.Errorf("could not get token from persistence: %w", err)
	}
	var cred login.Credential
	err = get.DataTo(&cred)
	if err != nil {
		return "", fmt.Errorf("could not parse data from persistence: %w", err)
	}
	return cred.User, nil
}

func (p *fireStorePersistence) DeleteToken(tk login.Token) error {
	collection := p.storage.Collection(loginTokensCollection)
	_, err := collection.Doc(string(tk)).Delete(p.ctx)
	if err != nil {
		return err
	}
	return nil
}

func (p *fireStorePersistence) SignUpUser(u login.UserName, pass login.Password) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pass), 12)
	if err != nil {
		return fmt.Errorf("could hash user password: %w", err)
	}
	collection := p.storage.Collection(loginCollection)
	_, err = collection.Doc(string(u)).Set(p.ctx, login.User{
		UserName: u,
		Password: login.Password(hashedPassword),
	})
	if err != nil {
		return fmt.Errorf("could not signup user: %w", err)
	}

	return nil
}

func (p *fireStorePersistence) ValidateCredentials(u login.UserName, pass login.Password) bool {
	collection := p.storage.Collection(loginCollection)
	var user login.User
	snapshot, err := collection.Doc(string(u)).Get(p.ctx)
	if err != nil {
		return false
	}
	err = snapshot.DataTo(&user)
	if err != nil {
		return false
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass))
	if err != nil {
		return false
	}
	return true
}

func (p *fireStorePersistence) DeleteUser(name login.UserName) error {
	collection := p.storage.Collection(loginCollection)
	_, err := collection.Doc(string(name)).Delete(p.ctx)
	if err != nil {
		return err
	}
	return nil
}

func (p *fireStorePersistence) GetProducts() (product.Products, error) {
	return nil, nil
}

func (p *fireStorePersistence) SaveProduct(product product.Product) error {
	//TODO implement me
	panic("implement me")
}

func (p *fireStorePersistence) DeleteProduct(id product.Id) error {
	//TODO implement me
	panic("implement me")
}

const loginCollection = "login"
const loginTokensCollection = "login-tokens"
