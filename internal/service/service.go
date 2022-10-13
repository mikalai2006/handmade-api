package service

import (
	"time"

	"github.com/mikalai2006/handmade/internal/domain"
	"github.com/mikalai2006/handmade/internal/repository"
	"github.com/mikalai2006/handmade/internal/utils"
	"github.com/mikalai2006/handmade/pkg/auths"
	"github.com/mikalai2006/handmade/pkg/hasher"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type Authorization interface {
	CreateAuth(auth domain.SignInInput) (primitive.ObjectID, error)
	SignIn(input domain.SignInInput) (domain.ResponseTokens, error)
	ExistAuth(auth domain.SignInInput) (domain.Auth, error)
	CreateSession(auth domain.Auth) (domain.ResponseTokens, error)
	VerificationCode(userId string, code string) error
	RefreshTokens(refreshToken string) (domain.ResponseTokens, error)
}

type Shop interface {
	Find(params domain.RequestParams) (domain.Response[domain.Shop], error)

	GetAllShops(params domain.RequestParams) (domain.Response[domain.Shop], error)
	CreateShop(userId string, shop domain.Shop) (*domain.Shop, error)
}

type User interface {
	GetUser(id string) (domain.User, error)
	FindUser(params domain.RequestParams) (domain.Response[domain.User], error)
	CreateUser(userId string, user domain.User) (*domain.User, error)
	DeleteUser(id string) (domain.User, error)
	UpdateUser(id string, user domain.User) (domain.User, error)
}


type Services struct {
	Authorization
	Shop
	User
}

type ConfigServices struct {
	Repositories *repository.Repositories
	Hasher hasher.PasswordHasher
	TokenManager auths.TokenManager
	OtpGenerator utils.Generator
	AccessTokenTTL time.Duration
	RefreshTokenTTL time.Duration
	VerificationCodeLength int
}

func NewServices(cfgService *ConfigServices) *Services {
	return &Services{
		Authorization: NewAuthService(cfgService.Repositories.Authorization, cfgService.Hasher, cfgService.TokenManager, cfgService.RefreshTokenTTL, cfgService.AccessTokenTTL, cfgService.OtpGenerator, cfgService.VerificationCodeLength),
		Shop: NewShopService(cfgService.Repositories.Shop),
		User: NewUserService(cfgService.Repositories.User),
	}
}