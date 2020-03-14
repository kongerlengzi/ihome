package handler

import (
	"08_ihome/apiv1.0/model"
	"08_ihome/apiv1.0/pkg/e"
	"08_ihome/apiv1.0/pkg/zaplog"
	"encoding/json"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"net/http"
	"strconv"
)

type House struct {
	Title string        `valid:"Required"`
	Price int           `valid:"Required"`
	Address string      `valid:"Required"`
	AreaID uint         `valid:"Required"`
}

func AddHouse(c *gin.Context)  {
	title := c.PostForm("title")
	price := c.PostForm("price")
	address := c.PostForm("address")
	areaid := c.PostForm("area_id")
	session, err := store.Get(c.Request,"login")
	if err != nil {
		log.Print(err)
		zaplog.Logger.Error(
			"sessiong get failed..",
			zap.String("session err","sessiong get failed"),
			zap.Error(err))
	}
	mobile := session.Values["mobile"].(string)
	u := model.GetUser(mobile)
	user_id := u.ID

	p,_ := strconv.Atoi(price)
	area_i,_ := strconv.ParseUint(areaid,10,8)
    area_id := uint(area_i)

	h := House{Title:title, Price:p, Address:address,AreaID:area_id}
	valid := validation.Validation{}
	b, err := valid.Valid(&h)
	if err != nil {
		log.Print(err)
	}
	if !b {
		for _, err := range valid.Errors {
			log.Println(err.Key, err.Message)
		}
		code := e.INVALID_PARAMS
		c.JSON(http.StatusOK,gin.H{
			"code" : code,
			"msg" : e.GetMsg(code),
		})
		return
	}
	data := make(map[string]interface{})
	data["user_id"] = user_id
	data["title"] = title
	data["price"] = p
	data["address"] = address
	data["area_id"] = area_id
	ok := model.AddHouse(data)
	if ok {
		code := e.SUCCESS
		c.JSON(http.StatusOK,gin.H{
			"code" : code,
			"msg" : e.GetMsg(code),
		})
	}
}

func AddArea(c *gin.Context)  {
	name := c.PostForm("name")
	data := make(map[string]interface{})
	data["name"] = name
	ok := model.AddArea(data)
	if ok {
		code := e.SUCCESS
		c.JSON(http.StatusOK,gin.H{
			"code" : code,
			"msg" : e.GetMsg(code),
		})
	}
}

func GetAllAreas(c *gin.Context)  {
	value := model.RedeisGet("area_info")
	if value["code"] == nil {
		areas := model.GetAllArea()
		code := e.SUCCESS
		value["code"] = code
		value["msg"] = e.GetMsg(code)
		value["data"] = areas
		value_r,_ := json.Marshal(value)
		model.RedisSet("100", value_r)
	}
	c.JSON(http.StatusOK, value)
}

func GetAllAreaHouses(c *gin.Context)  {
	id := c.PostForm("area_id")
	area_i,_ := strconv.ParseUint(id,10,8)
	area_id := uint(area_i)
	area, _ := model.GetArea(area_id)
	houses := model.GetAllAreaHouses(area)
	code := e.SUCCESS
	c.JSON(http.StatusOK,gin.H{
		"code" : code,
		"msg" : houses,
	})
}

func GetUserHouse(c *gin.Context)  {
	session, err := store.Get(c.Request,"login")
	if err != nil {
		log.Print(err)
	}
	mobile := session.Values["mobile"].(string)
	user := model.GetUser(mobile)
	house := model.GetUserHouse(user.ID)
	code := e.SUCCESS
	c.JSON(http.StatusOK,gin.H{
		"code" : code,
		"msg" : e.GetMsg(code),
		"data" : house,
	})
}

func IndexHouse(c *gin.Context)  {
	house := model.IndexHouse()
	code := e.SUCCESS
	c.JSON(http.StatusOK,gin.H{
		"code" : code,
		"msg" : e.GetMsg(code),
		"data" : house,
	})
}

func GetHouse(c *gin.Context)  {
	id := c.Param("house_id")
	idd,_ := strconv.Atoi(id)
	house_id := uint(idd)
	house := model.GetHouse(house_id)

	if house.ID <1 {
		code :=e.ERROR_EXIST_TAG
		c.JSON(http.StatusOK,gin.H{
			"code" : code,
			"msg" : e.GetMsg(code),
			"data" : "",
		})
	} else {
		code := e.SUCCESS
		c.JSON(http.StatusOK,gin.H{
			"code" : code,
			"msg" : e.GetMsg(code),
			"data" : house,
		})
	}
}