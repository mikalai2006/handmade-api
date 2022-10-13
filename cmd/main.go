package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/mikalai2006/handmade/internal/config"
	"github.com/mikalai2006/handmade/internal/handler"
	"github.com/mikalai2006/handmade/internal/repository"
	"github.com/mikalai2006/handmade/internal/server"
	"github.com/mikalai2006/handmade/internal/service"
	"github.com/mikalai2006/handmade/internal/utils"
	"github.com/mikalai2006/handmade/pkg/auths"
	"github.com/mikalai2006/handmade/pkg/hasher"
	"github.com/mikalai2006/handmade/pkg/logger"
	"github.com/sirupsen/logrus"
)

// @title Handmade API
// @version 1.0
// @description API Server for Handmade App

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	// setting logrus
	logrus.SetFormatter(new(logrus.JSONFormatter))

	// read config file
	cfg, err := config.Init("configs")
	if err != nil {
		logger.Error(err)
		return
	}

	// initialize mongoDB
	mongoClient, err := repository.NewMongoDB(repository.ConfigMongoDB{
		Host: cfg.Mongo.Host,
		Port: cfg.Mongo.Port,
		DBName: cfg.Mongo.Dbname,
		Username: cfg.Mongo.User,
		SSL: cfg.Mongo.SslMode,
		Password: cfg.Mongo.Password,
	})

	if err != nil {
		logger.Error(err)
	}

	mongoDB := mongoClient.Database(cfg.Mongo.Dbname)

	if (cfg.Environment != "prod") {
		logger.Info(mongoDB)
	}

	// initialize hasher
	hasher := hasher.NewSHA1Hasher(cfg.Auth.Salt)

	// initialize token manager
	tokenManager, err := auths.NewManager(cfg.Auth.SigningKey)
	if err != nil {
		logger.Error(err)

		return
	}

	// intiale opt
	otpGenerator := utils.NewGOTPGenerator()

	repositories := repository.NewRepositories(mongoDB)
	services := service.NewServices(&service.ConfigServices{
		Repositories: repositories,
		Hasher: hasher,
		TokenManager: tokenManager,
		OtpGenerator: otpGenerator,
		AccessTokenTTL: cfg.Auth.AccessTokenTTL,
		RefreshTokenTTL: cfg.Auth.RefreshTokenTTL,
		VerificationCodeLength: cfg.Auth.VerificationCodeLength,
	})
	handlers := handler.NewHandler(services, cfg.Oauth)

	// initialize server
	srv := server.NewServer(cfg, handlers.InitRoutes(*cfg))

	go func ()  {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("Error starting server: %s", err.Error())
		}
	}()

	logger.Infof("API service start on port: %s", cfg.HTTP.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<- quit

	logger.Info("API service shutdown")
	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		logger.Errorf("failed to stop server: %v", err)
	}

	if err := mongoClient.Disconnect(context.Background()); err != nil {
		logger.Error(err.Error())
	}
}
