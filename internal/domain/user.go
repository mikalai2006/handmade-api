package domain

import (
	"fmt"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Uid primitive.ObjectID `json:"uid,omitempty" bson:"uid,omitempty" db:"uid"`
	Type string `json:"type" db:"type" bson:"type"`
	Name string `json:"name,omitempty" db:"name" bson:"name" binding:"required"`
	Login string `json:"login" db:"login" bson:"login" binding:"required"`
	Currency string `json:"currency" bson:"currency" db:"currency"`
	Lang string `json:"lang" db:"lang" bson:"lang"`
	Online *bool `json:"online" db:"online" bson:"online"`
	Verify *bool `json:"verify" db:"verify" bson:"verify"`
	LastTime time.Time `json:"last_time" db:"last_time" bson:"last_time" form:"last_time"`
	CreatedAt time.Time `json:"created_at" db:"created_at" bson:"created_at" form:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" bson:"updated_at" form:"updated_at"`
}

func (user User) BodyToData() (interface{}, error) {

	result := bson.M{}
	var tagValue string
	elementsFilter := reflect.ValueOf(user)

	for i := 0; i < elementsFilter.NumField(); i += 1 {
		typeField := elementsFilter.Type().Field(i)
		tag := typeField.Tag

		tagValue = tag.Get("bson")

		if tagValue == "-" {
				continue
		}

		if elementsFilter.Field(i).Interface() == "" {
			continue
		}

		switch elementsFilter.Field(i).Kind() {
		case reflect.String:
				value := elementsFilter.Field(i).String()
				result[tagValue] = value

		case reflect.Bool:
				value := elementsFilter.Field(i).Bool()
				result[tagValue] = value

		case reflect.Int:
				value := elementsFilter.Field(i).Int()
				result[tagValue] = value

		}
	}
	if user.Online != nil {
		result["online"] = *user.Online
	}
	if user.Verify != nil {
		result["verify"] = *user.Verify
	}

	result["updated_at"] = time.Now()

	fmt.Println("user: new data =", result)

	return result, nil
}