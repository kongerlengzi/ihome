package router

import (
	"08_ihome/apiv1.0/handler"
	"08_ihome/apiv1.0/middleware/login_requied"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	gin.SetMode("debug")

	u := r.Group("/user")
	{
		u.POST("", handler.RegisterUser) //增加用户信息
		u.GET("/orders",login_requied.LoginRequied(),handler.GetUserOrders) //获取用户订单
		u.GET("/houses",login_requied.LoginRequied(),handler.GetUserHouse)  //获取用户房屋

	}

	a := r.Group("/area")
	{
		a.POST("",handler.AddArea) //增加地区信息
		a.POST("/house",handler.GetAllAreaHouses) //获取地区的房屋信息
	}

	r.POST("/login", handler.LoginUser)
	r.GET("/session", handler.GetSession)
	r.DELETE("/logout",login_requied.LoginRequied(), handler.LogoutUser)

    r.GET("/areas",handler.GetAllAreas) //获取所有地区信息
	r.GET("/houses/index",handler.IndexHouse)


	h := r.Group("/house")
	{
		h.POST("",login_requied.LoginRequied(),handler.AddHouse) //新增房屋信息
	    h.GET("/:house_id",handler.GetHouse) //获取单个房屋信息
	//	h.POST("/image",)
	//	h.GET("/search",)
	}

	o := r.Group("/orders",login_requied.LoginRequied())
	{
		o.POST("",handler.AddOrder)
		o.PUT("/status/:order_id",handler.UpdateOrderStatus)
		o.POST("/comment/:order_id",handler.SaveOrderComment)
	}

	return r
}