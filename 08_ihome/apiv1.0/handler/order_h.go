package handler

import (
	"08_ihome/apiv1.0/model"
	"08_ihome/apiv1.0/pkg/e"
	"08_ihome/apiv1.0/pkg/utils"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Order struct {
	UserID uint              `valid:"Required"`
	HousesID uint            `valid:"Required"`
	Begin_date time.Time     `valid:"Required"`
	End_date time.Time       `valid:"Required"`
}

func AddOrder(c *gin.Context)  {
	house_idd,_ := strconv.Atoi(c.PostForm("house_id"))
	house_id := uint(house_idd)
	start_dates := c.PostForm("start_date")
	end_dates := c.PostForm("end_date")
	start_date,err := time.Parse("2006-01-02",start_dates)
	end_date,err := time.Parse("2006-01-02",end_dates)

	user,_ := utils.GetSessionUser(c)
	user_id := user.ID

	order := Order{UserID:user_id,HousesID:house_id,Begin_date:start_date,End_date:end_date}
	valid := validation.Validation{}
	b, err := valid.Valid(&order)
	if err != nil {
		log.Print(err)
	}
	d := end_date.Sub(start_date)
	days := int(d.Hours()/24)

	if err != nil {
		code := e.INVALID_PARAMS
		c.JSON(http.StatusOK,gin.H{
			"code" : code,
			"msg" : "时间格式错误!",
		})
		return
	}
	if !b {
		code := e.INVALID_PARAMS
		c.JSON(http.StatusOK,gin.H{
			"code" : code,
			"msg" : e.GetMsg(code),
		})
		return
	}
	house := model.GetHouse(house_id)
	if house.ID <= 0 {
		code := e.INVALID_PARAMS
		c.JSON(http.StatusOK,gin.H{
			"code" : code,
			"msg" : "房屋不存在",
		})
		return
	}
	if house.UserID == user_id {
		code := e.INVALID_PARAMS
		c.JSON(http.StatusOK,gin.H{
			"code" : code,
			"msg" : "不能预定自己的房子",
		})
		return
	}
	data := make(map[string]interface{})
	amount := days * house.Price
	house_price := house.Price
	data["user_id"] = user_id
	data["house_id"] = house_id
	data["start_date"] = start_date
	data["end_date"] = end_date
	data["amount"] = amount
	data["days"] = days
	data["house_price"] = house_price

	conflictOrder := model.GetConflictOrder(data)
	if conflictOrder.ID > 0 {
		code := e.INVALID_PARAMS
		c.JSON(http.StatusOK,gin.H{
			"code" : code,
			"msg" : "房屋已经被预定",
		})
		return
	}
	orders := model.AddOrders(data)
	if orders {
		code := e.SUCCESS
		c.JSON(http.StatusOK,gin.H{
			"code" : code,
			"msg" : e.GetMsg(code),
		})
	} else {
		code := e.INVALID_PARAMS
		c.JSON(http.StatusOK,gin.H{
			"code" : code,
			"msg" : "订单添加失败",
		})
	}
}

func GetUserOrders(c *gin.Context)  {
	user,err := utils.GetSessionUser(c)
	if err != nil {
		code := e.INVALID_PARAMS
		c.JSON(http.StatusOK,gin.H{
			"code": code,
			"msg" : e.GetMsg(code),
		})
		return
	}
	user_id := user.ID

	role := c.DefaultQuery("role","")
	if role == "landlord" {
		houses := model.GetUserHouse(user_id)
		houses_id := []uint{}
		for _,house := range houses {
			houses_id = append(houses_id,house.ID )
		}
		orders := model.GetOrders(houses_id)
		code := e.SUCCESS
		c.JSON(http.StatusOK,gin.H{
			"code": code,
			"msg" : e.GetMsg(code),
			"orders":orders,
		})
	} else {
		orders := model.GetUserOrders(user_id)
		code := e.SUCCESS
		c.JSON(http.StatusOK,gin.H{
			"code": code,
			"msg" : e.GetMsg(code),
			"orders":orders,
		})
	}
}

func UpdateOrderStatus(c *gin.Context)  {
	user,_:= utils.GetSessionUser(c)
	user_id := user.ID

	order_d,_ := strconv.Atoi(c.Param("order_id"))
	action := c.PostForm("action")
	reason := c.PostForm("reason")
    if action != "accept" && action != "reject" {
		code := e.INVALID_PARAMS
		c.JSON(http.StatusOK,gin.H{
			"code" : code ,
			"msg" : e.GetMsg(code),
		})
		return
	} else {
		order_id := uint(order_d)
		status := "WAIT_ACCEPT"
		order := model.GetOrder(order_id,status)
		if order.ID < 1 || order.ID != user_id {
			code := e.ERROR
			c.JSON(http.StatusOK,gin.H{
				"code" : code ,
				"msg" : "订单不存在！",
			})
			return
		}
		if action == "accept" {
			status = "WAIT_PAYMENT"
			reason = ""
			ok := model.UpdateOrderStatus(status,reason,order_id)
			if ok {
				code := e.SUCCESS
				c.JSON(http.StatusOK,gin.H{
					"code" : code ,
					"msg" : e.GetMsg(code),
				})
			} else {
				code := e.ERROR
				c.JSON(http.StatusOK,gin.H{
					"code" : code ,
					"msg" : "failed！",
				})
			}
		}
		if action == "reject" {
			if reason == "" {
				code := e.ERROR
				c.JSON(http.StatusOK,gin.H{
					"code" : code ,
					"msg" : "resaon 不能为空！",
				})
				return
			}
			status = "REJECTED"
			ok := model.UpdateOrderStatus(status,reason,order_id)
			if ok {
				code := e.SUCCESS
				c.JSON(http.StatusOK,gin.H{
					"code" : code ,
					"msg" : e.GetMsg(code),
				})
			} else {
				code := e.ERROR
				c.JSON(http.StatusOK,gin.H{
					"code" : code ,
					"msg" : "failed！",
				})
			}
		}
	}
}

func SaveOrderComment(c *gin.Context)  {
	session, err := store.Get(c.Request,"login")
	if err != nil {
		log.Print(err)
	}
	mobile := session.Values["mobile"].(string)
	u := model.GetUser(mobile)
	user_id := u.ID

	order_d,_ := strconv.Atoi(c.Param("order_id"))
	order_id := uint(order_d)
	comment := c.PostForm("comment")
	if comment == "" {
		code := e.INVALID_PARAMS
		c.JSON(http.StatusOK,gin.H{
			"code" : code,
			"msg" : e.GetMsg(code),
		})
		return
	}
	status := "WAIT_COMMENT"

	orderC := model.GetOrderC(order_id, status, user_id)
    if orderC.ID <1 {
		code := e.INVALID_PARAMS
		c.JSON(http.StatusOK,gin.H{
			"code" : code,
			"msg" : e.GetMsg(code),
		})
		return
	}

	model.UpdateOrderStatus(status, comment, order_id)
	model.UpdateHouseOderById(orderC.HousesID)

	code := e.SUCCESS
	c.JSON(http.StatusOK,gin.H{
		"code" : code,
		"msg" : e.GetMsg(code),
	})
}