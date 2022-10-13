package service

import (
	"time"

	"github.com/mikalai2006/handmade/internal/domain"
	"github.com/mikalai2006/handmade/internal/repository"
	"github.com/mikalai2006/handmade/internal/utils"
	"github.com/mikalai2006/handmade/pkg/auths"
	"github.com/mikalai2006/handmade/pkg/hasher"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)


type AuthService struct {
	hasher hasher.PasswordHasher
	tokenManager auths.TokenManager

	repository repository.Authorization
	otpGenerator  utils.Generator

	accessTokenTTL         time.Duration
	refreshTokenTTL        time.Duration

	verificationCodeLength int
}

func NewAuthService(repo repository.Authorization, hasher hasher.PasswordHasher, tokenManager auths.TokenManager, refreshTokenTTL time.Duration, accessTokenTTL time.Duration, otp utils.Generator, verificationCodeLength int) *AuthService  {
	return &AuthService{
		hasher: hasher,
		tokenManager: tokenManager,

		repository: repo,
		otpGenerator: otp,

		accessTokenTTL: accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,

		verificationCodeLength: verificationCodeLength,
	}
}

func (s *AuthService) CreateAuth(auth domain.SignInInput) (primitive.ObjectID, error) {

	passwordHash, err := s.hasher.Hash(auth.Password)
	if err != nil {
		return primitive.NewObjectID(), err
	}

	verificationCode := s.otpGenerator.RandomSecret(s.verificationCodeLength)

	authData := domain.Auth{
		Login: auth.Login,
		Password: passwordHash,
		Email: auth.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Strategy: "local",
		Verification: domain.Verification{
			Code: verificationCode,
		},
	}

	 id, err := s.repository.CreateAuth(authData)
	 if err != nil {
		return primitive.NewObjectID(), err
	 }

	 // if created auth, send email with verification code

	 return id, nil
}


func (s *AuthService) ExistAuth(auth domain.SignInInput) (domain.Auth, error) {
	return s.repository.CheckExistAuth(auth)
}

func (s *AuthService) SignIn(auth domain.SignInInput) (domain.ResponseTokens, error)  {
	var result domain.ResponseTokens
	passwordHash, err := s.hasher.Hash(auth.Password)
	if err != nil {
		return result, err
	}
	auth.Password = passwordHash

	user, err := s.repository.GetByCredentials(auth)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return result, err
		}
		return result, err
	}

	return s.CreateSession(user)
}


func (s *AuthService) CreateSession(auth domain.Auth) (domain.ResponseTokens, error)  {
	var (
		res domain.ResponseTokens
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


func (s *AuthService) VerificationCode(userId string, hash string) error {
	err := s.repository.VerificationCode(userId, hash)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) RefreshTokens(refreshToken string) (domain.ResponseTokens, error) {
	var result domain.ResponseTokens

	user, err := s.repository.RefreshToken(refreshToken)
	if err != nil {
		return result, err
	}

	return s.CreateSession(user)
}