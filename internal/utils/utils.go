package utils

import (
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/handmade/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
)

type Utils interface {
	BodyToData() (interface{}, error)
	// ParamsToFilter() (interface{}, error)
}

// Create interface from body request to update item mongodb
func GetBodyToData(u Utils) (interface{}, error) {
	data, err := u.BodyToData()
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Parse request params and return struct domain.RequestParams
func GetParamsFromRequest(c *gin.Context, filterStruct interface{}) (domain.RequestParams, error) {
	params := domain.RequestParams{
		Filter: filterStruct,
	}
	var filter domain.Shop
	if err := c.Bind(&filter); err != nil {
		return domain.RequestParams{}, err
	}
	// filterBson, err := bson.Marshal(filter)
	// var interfaceFilter interface{}
	// bson.Unmarshal(filterBson, &interfaceFilter)
	// for k,v := range interfaceFilter {
	// }
	dataFilter := bson.M{}
	var tagValue string
	elementsFilter := reflect.ValueOf(filter)
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
			dataFilter[tagValue] = value

		case reflect.Bool:
			value := elementsFilter.Field(i).Bool()
			dataFilter[tagValue] = value

		case reflect.Int:
			value := elementsFilter.Field(i).Int()
			dataFilter[tagValue] = value
		}
	}

	var opts domain.Options
	if err := c.Bind(&opts); err != nil {
		return domain.RequestParams{}, err
	}

	sort := c.QueryMap("$sort")
	if len(sort) > 0 {
		var testBson bson.D
		for k, v := range sort {
			value, err := strconv.ParseInt(v, 10, 32)
			if err != nil {
				return domain.RequestParams{}, err
			}
			testBson = append(testBson, bson.E{Key: k, Value: value})
		}
		opts.Sort = testBson
	}

	// err = bson.Unmarshal(sort, &sort)
	// fmt.Println("----------")
	// fmt.Printf("dataFilter=%s", dataFilter)
	// fmt.Println("----------")
	// fmt.Printf("len dataFilter=%s", len(dataFilter))
	// fmt.Println("----------")
	// fmt.Printf("filter=%s", filter)
	// fmt.Println("----------")
	// fmt.Printf("sort=%s", testBson)
	// fmt.Println("----------")
	// fmt.Printf("opts=%s", opts)
	// fmt.Println("----------")

	if opts.Limit == 0 || opts.Limit > 50 {
		opts.Limit = 10
	}

	params.Filter = dataFilter
	params.Options = opts

	return params, nil
}