package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/mikalai2006/handmade/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthMongo struct {
	db *mongo.Database
}

func NewAuthMongo(db *mongo.Database) *AuthMongo {
	return &AuthMongo{db:db}
}

func (r *AuthMongo) CreateAuth(user domain.SignInInput) (primitive.ObjectID, error) {
	var id primitive.ObjectID

	collection := r.db.Collection(tblAuth)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	auth, err := collection.InsertOne(ctx, user)
	if err != nil {
		return id, err
	}
	id = auth.InsertedID.(primitive.ObjectID)

	return id, nil
}

func chooseProvider(auth domain.SignInInput) (bson.D) {
	if auth.Strategy == "local" {
		return bson.D{{Key: "login", Value: auth.Login}, {Key: "password", Value: auth.Password}}
	}
	if auth.VkId != "" {
		return bson.D{{Key: "vkid", Value: auth.VkId}}
	} else if auth.GoogleId != "" {
		return bson.D{{Key: "googleid", Value: auth.GoogleId}}
	}
	return bson.D{{Key: "vkid", Value: "none"}}
}

func (r *AuthMongo) CheckExistAuth(auth domain.SignInInput) (domain.Auth, error) {
	var user domain.Auth

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	filter := chooseProvider(auth)
	fmt.Println("")
	fmt.Printf("chooseProvider: filter=%s", filter)
	err := r.db.Collection(tblAuth).FindOne(ctx, filter).Decode(&user)
	if err != nil {
    // ErrNoDocuments means that the filter did not match any documents in the collection
    if err == mongo.ErrNoDocuments {
			err = nil
    }
}
	return user, err
}

func (r *AuthMongo) GetAuth(auth domain.Auth) (domain.Auth, error) {
	var user domain.Auth

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	query := bson.M{"login": auth.Login, "password": auth.Password}
	// fmt.Println("")
	// fmt.Printf("GetAuth: query=%s", query)
	err := r.db.Collection(tblAuth).FindOne(ctx, query).Decode(&user)
	// if err != nil {
	// 	return domain.Auth{}, err
	// }

	return user, err
}

func (r *AuthMongo) GetByCredentials(auth domain.SignInInput) (domain.Auth, error) {
	var user domain.Auth

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	// query := bson.M{"login": auth.Login, "password": auth.Password}
	filter := chooseProvider(auth)
	// fmt.Println("---")
	// fmt.Println(filter)
	err := r.db.Collection(tblAuth).FindOne(ctx, filter).Decode(&user)

	return user, err
}

func (r *AuthMongo) SetSession(authID primitive.ObjectID, session domain.Session) error {

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	_, err := r.db.Collection(tblAuth).UpdateOne(ctx, bson.M{"_id": authID}, bson.M{"$set": bson.M{"session": session, "lastVisitAt": time.Now()}})

	return err
}
