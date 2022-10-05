package repository

import (
	"context"

	"github.com/mikalai2006/handmade/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ShopMongo struct {
	db *mongo.Database
}

const (
	collectionName string = "shops"
)


func NewShopMongo(db *mongo.Database) *ShopMongo {
	return &ShopMongo{db:db}
}


func (r *ShopMongo) GetAllShops() ([]*domain.Shop, error) {
	var results []*domain.Shop

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	cursor, err := r.db.Collection(collectionName).Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)


	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *ShopMongo) CreateShop(userId string, shop domain.Shop) (*domain.Shop, error) {
	var result *domain.Shop

	collection := r.db.Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	newShop := domain.Shop{
		Title: shop.Title,
		Description: shop.Description,
		Seo: "",
		UserId: userId,
	}

	res, err := collection.InsertOne(ctx, newShop)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(collectionName).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}