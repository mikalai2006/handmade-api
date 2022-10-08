package handler

import (
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/handmade/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
)

func (h *Handler) CreateShop(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	var input domain.Shop
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	shop, err := h.services.Shop.CreateShop(userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, shop)
}

// @Summary Shop Get all shops
// @Security ApiKeyAuth
// @Tags shop
// @Description get all shops
// @ModuleID shops
// @Accept  json
// @Produce  json
// @Success 200 {object} []domain.Shop
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/shops [get]
func (h *Handler) GetAllShops(c *gin.Context) {
	var params domain.RequestParams
	var filter domain.Shop
	if err := c.Bind(&filter); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
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
		newErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	sort := c.QueryMap("$sort")
	var testBson bson.D
	for k,v := range sort {
		value, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())

			return
		}
		testBson = append(testBson, bson.E{Key: k, Value: value})
	}
	// err = bson.Unmarshal(sort, &sort)
	params.Filter = dataFilter
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

	if opts.Limit == 0 {
		opts.Limit = 10
	}
	opts.Sort = testBson
	params.Options = opts

	shops, err := h.services.Shop.GetAllShops(params)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, shops)
}

func (h *Handler) GetShopById(c *gin.Context) {

}

func (h *Handler) UpdateShop(c *gin.Context) {

}

func (h *Handler) DeleteShop(c *gin.Context) {

}