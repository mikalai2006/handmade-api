package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	tblShops = "shops"
	tblUsers = "users"
	tblAuth = "auth"
	MongoQueryTimeout = 10 * time.Second
)


type ConfigMongoDB struct {
	Host string
	Port string
	Username string
	Password string
	DBName string
	// SSL string
}



func NewMongoDB(cfg ConfigMongoDB) (*mongo.Database, error) {
	// ctx will be used to set deadline for process, here
	// deadline will of 30 seconds.
	ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
	// Release resource when the main
	// function is returned.
	defer cancel()

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=admin&readPreference=primary&directConnection=true&ssl=false",
	cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	logrus.Println(uri)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
		if err != nil {
			return nil, err
		}


		if err := client.Ping(ctx, readpref.Primary()); err != nil {
			return nil, err
		}


	return client.Database(cfg.DBName), nil
}
