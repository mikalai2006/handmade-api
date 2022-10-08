package domain

import "go.mongodb.org/mongo-driver/bson"

type Response struct {
	Total int64 `json:"total" bson:"total"`
	Limit int   `json:"limit" bson:"limit"`
	Skip  int   `json:"skip" bson:"skip"`
	Data  []any `json:"data" bson:"data"`
}

type RequestParams struct {
	Options
	Filter interface{} `json:"filter" bson:"filter"`
	Group interface{} `json:"group" bson:"$group"`
}

type Options struct {
	Limit int64  `json:"limit" bson:"limit" form:"$limit"`
	Skip  int64  `json:"skip" bson:"skip" form:"$skip"`
	Sort  bson.D `json:"sort" bson:"sort" form:"$sort"`
}
