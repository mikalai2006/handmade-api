package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/mikalai2006/handmade/internal/domain"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
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
	SSL bool
}



func NewMongoDB(cfg ConfigMongoDB) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
	defer cancel()

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=admin&readPreference=primary&directConnection=true&ssl=%t",
	cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSL)
	logrus.Println(uri)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
		if err != nil {
			return nil, err
		}

		if err := client.Ping(ctx, readpref.Primary()); err != nil {
			return nil, err
		}


	return client, nil
}


func CreatePipeline(params domain.RequestParams) (mongo.Pipeline, error) {
	pipe := mongo.Pipeline{}
	pipe = append(pipe, bson.D{{Key: "$match", Value: params.Filter}})
	// opts := options.Find()
	if params.Options.Sort != nil {
		// opts.SetSort(params.Options.Sort)
		pipe = append(pipe, bson.D{{Key: "$sort", Value: params.Options.Sort}})
	}
	if params.Options.Skip != 0 {
		// opts.SetSkip(params.Options.Skip)
		pipe = append(pipe, bson.D{{Key: "$skip", Value: params.Options.Skip}})
	}
	if params.Options.Limit != 0 {
		// opts.SetLimit(params.Options.Limit)
		pipe = append(pipe, bson.D{{Key: "$limit", Value: params.Options.Limit}})
	}

	// pipe = append(pipe, bson.D{
	// 	{Key: "$group", Value: bson.M{
	// 		"_id":    "$title",
	// 		"count": bson.M{"$sum": 1},
	// }}})

	return pipe, nil
}

func CreateFilterAndOptions(params domain.RequestParams) (interface{}, *options.FindOptions, error) {
	opts := options.Find()
	if params.Options.Sort != nil {
		opts.SetSort(params.Options.Sort)
	}
	if params.Options.Skip != 0 {
		opts.SetSkip(params.Options.Skip)
	}
	if params.Options.Limit != 0 {
		opts.SetLimit(params.Options.Limit)
	}

	filter := params.Filter

	return filter, opts, nil
}