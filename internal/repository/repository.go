package repository

import (
	"github.com/mikalai2006/handmade/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Authorization interface {
	CreateAuth(auth domain.Auth) (primitive.ObjectID, error)
	GetAuth(auth domain.Auth) (domain.Auth, error)
	CheckExistAuth(auth domain.Auth) (domain.Auth, error)
	GetByCredentials(auth domain.Auth) (domain.Auth, error)
	SetSession(authId primitive.ObjectID, session domain.Session) error
}

type Shop interface {
	GetAllShops() ([]*domain.Shop, error)
	CreateShop(userId string, shop domain.Shop) (*domain.Shop, error)
}


type TodoList interface {
}

type TodoItem interface {
}


type Repository struct {
	Authorization
	TodoList
	TodoItem
	Shop
}

func NewRepository(mongodb *mongo.Database) *Repository {
	return &Repository{
		Authorization: NewAuthMongo(mongodb),
		Shop: NewShopMongo(mongodb),
	}
}