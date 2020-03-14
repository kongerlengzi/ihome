package handler

import (
	"08_ihome/apiv1.0/model"
	"08_ihome/apiv1.0/pkg/e"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

type Register struct {
	Name string     `binding:"Required",json:"name",xml:"name",form:"name"`
	Mobile string   `binding:"Required",json:"mobile",xml:"mobile",form:"mobile"`
	Pwd string      `binding:"Required",json:"pwd",xml:"pwd",form:"pwd"`
}

type Login struct {
	Mobile string   `binding:"Required",json:"mobile",xml:"mobile",form:"mobile" `
	Pwd string      `binding:"Required",json:"pwd",xml:"pwd",form:"pwd"`
}

var store = sessions.NewCookieStore([]byte("something-very-secret"))

func RegisterUser(c *gin.Context)  {
	var r Register
	if err := c.ShouldBind(&r);err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data := make(map[string]interface{})
	IsExit := model.ExistUserByID(r.Mobile)
	if IsExit {
		code := e.ERROR_AUTH
		c.JSON(http.StatusOK,gin.H{
			"code" : code,
			"msg" : e.GetMsg(code),
			"data" : data,
		})
		return
	}
	ok := model.AddUser(r)
	if ok {
		code := e.SUCCESS
		c.JSON(http.StatusOK,gin.H{
			"code" : code,
			"msg" : e.GetMsg(code),
		})
	}
}

func LoginUser(c *gin.Context) {
	var l Login
	if err := c.ShouldBind(&l);err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	IsLogin := model.Login(l.Mobile,l.Pwd)
	if IsLogin {
		session,err := store.New(c.Request,"login")
		if err != nil {
			log.Print(err)
		}
		session.Values["mobile"] = l.Mobile
		session.Options.MaxAge = 0
		_ = session.Save(c.Request, c.Writer)

		code := e.SUCCESS
		c.JSON(http.StatusOK,gin.H{
			"code" : code,
			"msg" : e.GetMsg(code),
		})
	} else {
		code := e.ERROR_AUTH
		c.JSON(http.StatusOK,gin.H{
			"code" : code,
			"msg" : e.GetMsg(code),
		})
	}
}

func LogoutUser(c *gin.Context) {
	session, err := store.Get(c.Request,"login")
	if err != nil {
		log.Print(err)
	}
	session.Options.MaxAge = -1
	_ = session.Save(c.Request, c.Writer)

	code := e.SUCCESS
	c.JSON(http.StatusOK,gin.H{
		"code" : code,
		"msg" : e.GetMsg(code),
	})

}

func GetSession (c *gin.Context) {
	session, err := store.Get(c.Request,"login")
	if err != nil {
		log.Print(err)
	}
	mobile := session.Values["mobile"]
	c.JSON(http.StatusOK,gin.H{
		"code" : 200,
		"msg" : mobile,
	})
}