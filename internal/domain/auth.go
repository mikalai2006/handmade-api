package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Auth struct {
	// swagger:ignore
	Id       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Login    string    `json:"login" binding:"required"`
	Password string    `json:"password" binding:"required"`
	Strategy string    `json:"-"`
	VkId     string `json:"-"`
	GoogleId string    `json:"-"`
	GithubId string    `json:"-"`
	AppleId  string    `json:"-"`
	Verification     Verification         `json:"verification" bson:"verification"`
	Session          Session              `json:"session" bson:"session,omitempty"`
}

type Verification struct {
	Code     string `json:"code" bson:"code"`
	Verified bool   `json:"verified" bson:"verified"`
}