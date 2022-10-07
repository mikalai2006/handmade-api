package service

import (
	"time"

	"github.com/mikalai2006/handmade/internal/domain"
	"github.com/mikalai2006/handmade/internal/repository"
	"github.com/mikalai2006/handmade/pkg/auths"
	"github.com/mikalai2006/handmade/pkg/hasher"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type AuthService struct {
	hasher hasher.PasswordHasher
	tokenManager auths.TokenManager

	repository repository.Authorization

	accessTokenTTL         time.Duration
	refreshTokenTTL        time.Duration
}

func NewAuthService(repo repository.Authorization, hasher hasher.PasswordHasher, tokenManager auths.TokenManager, refreshTokenTTL time.Duration, accessTokenTTL time.Duration) *AuthService  {
	return &AuthService{
		hasher: hasher,
		tokenManager: tokenManager,

		repository: repo,

		accessTokenTTL: accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (s *AuthService) CreateAuth(auth domain.SignInInput) (primitive.ObjectID, error) {

	passwordHash, err := s.hasher.Hash(auth.Password)
	if err != nil {
		return primitive.NewObjectID(), err
	}
	auth.Password = passwordHash
	return s.repository.CreateAuth(auth)
}


func (s *AuthService) ExistAuth(auth domain.SignInInput) (domain.Auth, error) {
	return s.repository.CheckExistAuth(auth)
}

func (s *AuthService) SignIn(auth domain.SignInInput) (Tokens, error)  {
	passwordHash, err := s.hasher.Hash(auth.Password)
	if err != nil {
		return Tokens{}, err
	}
	auth.Password = passwordHash

	user, err := s.repository.GetByCredentials(auth)
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

	res.AccessToken, err = s.tokenManager.NewJWT(auth.Id.Hex(), s.accessTokenTTL)
	if err != nil {
		return res, err
	}

	res.RefreshToken, err = s.tokenManager.NewRefreshToken()
	if err != nil {
		return res, err
	}

	session := domain.Session{
		RefreshToken: res.RefreshToken,
		ExpiresAt:    time.Now().Add(s.refreshTokenTTL),
	}

	err = s.repository.SetSession(auth.Id, session)

	return res, err
}