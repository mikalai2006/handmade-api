package handler

import (
	"errors"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/handmade/internal/domain"
	"github.com/mikalai2006/handmade/pkg/auths"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	authorizationHeader = "Authorization"
	userCtx = "userId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)

	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	if len(headerParts[1]) == 0 {
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	// parse token
	// userId, err := h.services.Authorization.ParseToken(headerParts[1])
	// if err != nil {
	// 	newErrorResponse(c, http.StatusUnauthorized, err.Error())
	// 	return
	// }
	tokenManager, err := auths.NewManager(os.Getenv("SIGNING_KEY"))
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	id, err := tokenManager.Parse(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	c.Set(userCtx, id)
	// session := sessions.Default(c)
	// user := session.Get(userkey)
	// if user == nil {
	// 	// Abort the request with the appropriate error code
	// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	// 	return
	// }
	// logrus.Printf("user session= %s", user)
	// // Continue down the chain to handler etc
	// c.Next()
}

func getUserId(c *gin.Context) (string, error)  {
	id, ok := c.Get(userCtx)
	if !ok {
		newErrorResponse(c, http.StatusUnauthorized, "user id not found")
		return "", errors.New("user not found")
	}

	idInt, ok := id.(string)
	if !ok {
		newErrorResponse(c, http.StatusUnauthorized, "user id is of invalid type")
		return "", errors.New("user not found2")
	}

	return idInt, nil
}

// Parse request params and return struct domain.RequestParams
func getParamsFromRequest(c *gin.Context, filterStruct interface{}) (domain.RequestParams, error) {
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
	if (len(sort) > 0) {
		var testBson bson.D
		for k,v := range sort {
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