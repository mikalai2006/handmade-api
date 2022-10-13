package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Auth struct {
	// swagger:ignore
	Id       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Login    string    `json:"login" binding:"required"`
	Email 		string `json:"email"`
	Password string    `json:"password" binding:"required"`
	Strategy string    `json:"-"`
	VkId     string `json:"-"`
	GoogleId string    `json:"-" bson:"google_id"`
	GithubId string    `json:"-"`
	AppleId  string    `json:"-"`
	Verification     Verification         `json:"verification" bson:"verification"`
	Session          Session              `json:"session" bson:"session,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type Verification struct {
	Code     string `json:"code" bson:"code"`
	Verified bool   `json:"verified" bson:"verified"`
}
type SignInInput struct {
	Login string `json:"login" bson:"login"`
	Email string `json:"-" bson:"-"`
	Password string `json:"password" bson:"password"`
	Strategy string `json:"strategy" bson:"strategy"`
	VkId string `json:"-"`
	GoogleId string `json:"-"`
}