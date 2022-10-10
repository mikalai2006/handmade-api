package domain

import "time"

// Id          int    `json:"id" bson:"_id"`
type Shop struct {
	Title       string `json:"title" bson:"title" form:"title"`
	Description string `json:"description" bson:"description" form:"description"`
	Seo         string `json:"seo" bson:"seo" form:"seo"`
	UserId      string `json:"user_id" bson:"user_id" form:"user_id"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at" form:"created_at"`
}