package utils

import (
	"08_ihome/apiv1.0/model"
	"08_ihome/apiv1.0/pkg/zaplog"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
	"log"
)

var store = sessions.NewCookieStore([]byte("something-very-secret"))

func GetSessionUser (c *gin.Context) (user model.User,err error) {
	session, err := store.Get(c.Request,"login")
	if err != nil {
		log.Print(err)
	}
	if mobile := session.Values["mobile"] ; mobile == nil{
		err = errors.New("session fialed!")
		log.Print(err)
		zaplog.Logger.Error(
			"session not found",
			zap.String("session","fialed"))
		return  user,err
	} else {
		mobile := session.Values["mobile"].(string)
		user = model.GetUser(mobile)
		return user,nil
	}
}
