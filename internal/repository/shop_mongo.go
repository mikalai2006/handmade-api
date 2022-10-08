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


func (r *ShopMongo) GetAllShops(params domain.RequestParams) (domain.Response, error) {
	var results []domain.Shop
	var response domain.Response

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

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
	cursor, err := r.db.Collection(collectionName).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &results); err != nil {
		return response, err
	}

	var resultSlice []interface{} = make([]interface{}, len(results))
	for i, d := range results {
		resultSlice[i] = d
	}

	count,err := r.db.Collection(collectionName).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response{
		Total: count,
		Skip: int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data: resultSlice,
	}
	return response, nil
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