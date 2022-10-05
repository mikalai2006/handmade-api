package service

import (
	"github.com/mikalai2006/handmade/internal/domain"
	"github.com/mikalai2006/handmade/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type Authorization interface {
	CreateAuth(auth domain.Auth) (primitive.ObjectID, error)
	SignIn(auth domain.Auth) (Tokens, error)
	ExistAuth(auth domain.Auth) (domain.Auth, error)
	CreateSession(auth domain.Auth) (Tokens, error)
}

type Shop interface {
	GetAllShops() ([]*domain.Shop, error)
	CreateShop(userId string, shop domain.Shop) (*domain.Shop, error)
}

type TodoList interface {
}

type TodoItem interface {
}

type Service struct {
	Authorization
	TodoList
	TodoItem
	Shop
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Shop: NewShopService(repos.Shop),
	}
}