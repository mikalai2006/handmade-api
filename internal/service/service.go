package service

import (
	"time"

	"github.com/mikalai2006/handmade/internal/domain"
	"github.com/mikalai2006/handmade/internal/repository"
	"github.com/mikalai2006/handmade/pkg/auths"
	"github.com/mikalai2006/handmade/pkg/hasher"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type Authorization interface {
	CreateAuth(auth domain.SignInInput) (primitive.ObjectID, error)
	SignIn(input domain.SignInInput) (Tokens, error)
	ExistAuth(auth domain.SignInInput) (domain.Auth, error)
	CreateSession(auth domain.Auth) (Tokens, error)
}

type Shop interface {
	Find(params domain.RequestParams) (domain.Response, error)

	GetAllShops(params domain.RequestParams) (domain.Response, error)
	CreateShop(userId string, shop domain.Shop) (*domain.Shop, error)
}

type Services struct {
	Authorization
	Shop
}

type ConfigServices struct {
	Repositories *repository.Repositories
	Hasher hasher.PasswordHasher
	TokenManager auths.TokenManager
	AccessTokenTTL time.Duration
	RefreshTokenTTL time.Duration
}

func NewServices(cfgService *ConfigServices) *Services {
	return &Services{
		Authorization: NewAuthService(cfgService.Repositories.Authorization, cfgService.Hasher, cfgService.TokenManager, cfgService.RefreshTokenTTL, cfgService.AccessTokenTTL),
		Shop: NewShopService(cfgService.Repositories.Shop),
	}
}