package domain

type Shop struct {
	// Id          int    `json:"id" bson:"_id"`
	Title       string `json:"title" bson:"title" form:"title"`
	Description string `json:"description" bson:"description" form:"description"`
	Seo         string `json:"seo" bson:"seo" form:"seo"`
	UserId      string `json:"user_id" bson:"user_id" form:"user_id"`
}