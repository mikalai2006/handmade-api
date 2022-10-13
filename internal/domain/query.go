package domain

type Response[D any] struct {
	Total int64 `json:"total" bson:"total"`
	Limit int   `json:"limit" bson:"limit"`
	Skip  int   `json:"skip" bson:"skip"`
	Data  []D   `json:"data" bson:"data"`
}

type RequestParams struct {
	Options
	Filter interface{} `json:"filter" bson:"filter"`
	Group  interface{} `json:"group" bson:"$group"`
}

type Options struct {
	Limit int64       `json:"$limit" bson:"limit" form:"$limit"`
	Skip  int64       `json:"$skip" bson:"skip" form:"$skip"`
	Sort  interface{} `json:"$sort" bson:"sort" form:"$sort"`
}

type RefreshInput struct {
	Token string `json:"token" bson:"token" form:"token" binding:"required"`
}