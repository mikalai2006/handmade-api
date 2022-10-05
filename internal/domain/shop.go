package domain

type Shop struct {
	// Id          int    `json:"id" bson:"_id"`
	Title       string `json:"title" bson:"title" binding:"required"`
	Description string `json:"description" bson:"description"`
	Seo         string `json:"seo" bson:"seo"`
	UserId      string `json:"user_id" bson:"user_id"`
}