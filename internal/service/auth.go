package service

import (
	"os"
	"time"

	"github.com/mikalai2006/handmade/internal/domain"
	"github.com/mikalai2006/handmade/internal/repository"
	"github.com/mikalai2006/handmade/pkg/auths"
	"github.com/mikalai2006/handmade/pkg/hasher"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthService struct {
	repo repository.Authorization
	accessTokenTTL         time.Duration
	refreshTokenTTL        time.Duration
}

func NewAuthService(repo repository.Authorization) *AuthService  {
	return &AuthService{
		repo: repo,
		accessTokenTTL: viper.GetDuration("auth.accessTokenTTL"),
		refreshTokenTTL: viper.GetDuration("auth.refreshTokenTTL"),
	}
}

func (s *AuthService) CreateAuth(auth domain.Auth) (primitive.ObjectID, error) {

	// init salt
	hasher := hasher.NewSHA1Hasher(os.Getenv("SALT"))
	passwordHash, _ := hasher.Hash(auth.Password)
	// if err != nil {
	// 	return "0", err
	// }
	auth.Password = passwordHash
	return s.repo.CreateAuth(auth)
}


func (s *AuthService) ExistAuth(auth domain.Auth) (domain.Auth, error) {
	return s.repo.CheckExistAuth(auth)
}

func (s *AuthService) SignIn(auth domain.Auth) (Tokens, error)  {
	// init salt
	hasher := hasher.NewSHA1Hasher(viper.GetString("salt"))
	passwordHash, err := hasher.Hash(auth.Password)
	if err != nil {
		return Tokens{}, nil
	}
	auth.Password = passwordHash

	user, err := s.repo.GetByCredentials(auth)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return Tokens{}, err
		}
		return Tokens{}, err
	}

	return s.CreateSession(user)
}


func (s *AuthService) CreateSession(auth domain.Auth) (Tokens, error)  {
	var (
		res Tokens
		err error
	)

	tokenManager, err := auths.NewManager(os.Getenv("SIGNING_KEY"))
	if err != nil {
		return res, err
	}

	res.AccessToken, err = tokenManager.NewJWT(auth.Id.Hex(), s.accessTokenTTL)
	if err != nil {
		return res, err
	}

	res.RefreshToken, err = tokenManager.NewRefreshToken()
	if err != nil {
		return res, err
	}

	session := domain.Session{
		RefreshToken: res.RefreshToken,
		ExpiresAt:    time.Now().Add(s.refreshTokenTTL),
	}

	err = s.repo.SetSession(auth.Id, session)

	return res, err
}