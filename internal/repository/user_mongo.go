package repository

import (
	"context"
	"time"

	"github.com/mikalai2006/handmade/internal/domain"
	"github.com/mikalai2006/handmade/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserMongo struct {
	db *mongo.Database
}

const (
	userCollection string = "users"
)


func NewUserMongo(db *mongo.Database) *UserMongo {
	return &UserMongo{db:db}
}


func (r *UserMongo) GetUser(id string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result domain.User

	userIdPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.User{}, err
	}

	filter := bson.M{"_id": userIdPrimitive}

	err = r.db.Collection(userCollection).FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return domain.User{}, err
	}

	return result, nil
}

func (r *UserMongo) FindUser(params domain.RequestParams) (domain.Response[domain.User], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.User
	var response domain.Response[domain.User]
	pipe, err := CreatePipeline(params)
	if err != nil {
		return domain.Response[domain.User]{}, err
	}

	cursor, err := r.db.Collection(userCollection).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &results); err != nil {
		return response, err
	}

	var resultSlice []domain.User = make([]domain.User, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count,err := r.db.Collection(userCollection).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[domain.User]{
		Total: count,
		Skip: int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data: resultSlice,
	}
	return response, nil
}

func (r *UserMongo) CreateUser(userId string, user domain.User) (*domain.User, error) {
	var result *domain.User

	collection := r.db.Collection(userCollection)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIdPrimitive, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}

	newUser := domain.User{
		Name: user.Name,
		Uid: userIdPrimitive,
		Type: "guest",
		Login: user.Login,
		Lang: user.Lang,
		Currency: user.Currency,
		Online: user.Online,
		Verify: user.Verify,
		LastTime: time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	res, err := collection.InsertOne(ctx, newUser)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(userCollection).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}


func (r *UserMongo) DeleteUser(id string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = domain.User{}
	collection := r.db.Collection(userCollection)

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": idPrimitive}

	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}

	_, err = collection.DeleteOne(ctx, filter)
	if err != nil  {
		return result, err
	}

	return result, nil
}


func (r *UserMongo) UpdateUser(id string, user domain.User) (domain.User, error) {
	var result domain.User
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(userCollection)

	data, err := utils.GetBodyToData(user)
	if err != nil {
		return result, err
	}

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": idPrimitive}

	// fmt.Println("data=", data)
	_, err = collection.UpdateOne(ctx, filter, bson.M{"$set": data})
	if err != nil {
		return result, err
	}

	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}