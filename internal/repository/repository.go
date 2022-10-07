package repository

import (
	"github.com/mikalai2006/handmade/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Authorization interface {
	CreateAuth(auth domain.SignInInput) (primitive.ObjectID, error)
	GetAuth(auth domain.Auth) (domain.Auth, error)
	CheckExistAuth(auth domain.SignInInput) (domain.Auth, error)
	GetByCredentials(auth domain.SignInInput) (domain.Auth, error)
	SetSession(authId primitive.ObjectID, session domain.Session) error
}

type Shop interface {
	GetAllShops() (domain.Response, error)
	CreateShop(userId string, shop domain.Shop) (*domain.Shop, error)
}

type Repositories struct {
	Authorization
	Shop
}

func NewRepositories(mongodb *mongo.Database) *Repositories {
	return &Repositories{
		Authorization: NewAuthMongo(mongodb),
		Shop: NewShopMongo(mongodb),
	}
}