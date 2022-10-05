package main

import (
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mikalai2006/handmade"
	"github.com/mikalai2006/handmade/internal/handler"
	"github.com/mikalai2006/handmade/internal/repository"
	"github.com/mikalai2006/handmade/internal/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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
	if err := initConfig(); err != nil {
		logrus.Fatalf("Error initializing config: %s", err.Error())
	}

	// read env configs
	if err :=godotenv.Load(); err != nil {
		logrus.Fatalf("Error loading env envirables: %s", err.Error())
	}

	// initialize postgres
	// db, err := repository.NewPostgresDB(repository.Config{
	// 	Host: os.Getenv("PG_HOST"),
	// 	Port: os.Getenv("PG_PORT"),
	// 	DBName: viper.GetString("db.dbname"),
	// 	Username: os.Getenv("PG_USER"),
	// 	SSLMode: viper.GetString("db.sslmode"),
	// 	Password: os.Getenv("PG_PASSWORD"),
	// })
	// if err != nil {
	// 	logrus.Fatalf("Failed to initialize postgres db: %s", err.Error())
	// }

	// initialize mongoDB
	mongoDB, err := repository.NewMongoDB(repository.ConfigMongoDB{
		Host: os.Getenv("MONGODB_HOST"),
		Port: os.Getenv("MONGODB_PORT"),
		DBName: viper.GetString("mongodb.dbname"),
		Username: os.Getenv("MONGODB_USER"),
		// SSL: viper.GetString("mongodb.ssl"),
		Password: os.Getenv("MONGODB_PASSWORD"),
	})
	if err != nil {
		logrus.Fatalf("Failed to initialize mongoDB: %s", err.Error())
	}
	logrus.Debug(mongoDB)


	repos := repository.NewRepository(mongoDB)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(handmade.Server)
	if err := srv.Run(os.Getenv("PORT"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("Error starting server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}