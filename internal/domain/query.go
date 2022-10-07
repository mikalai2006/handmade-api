package domain

type Response struct {
	Total int64 `json:"total" bson:"total"`
	Limit int64 `json:"limit" bson:"limit"`
	Skip  int64 `json:"skip" bson:"skip"`
	Data  []any `json:"data" bson:"data"`
}

type Pagination struct {
	Limit  int64       `json:"limit" bson:"limit"`
	Skip   int64       `json:"skip" bson:"skip"`
	Filter interface{} `json:"filter" bson:"filter"`
}