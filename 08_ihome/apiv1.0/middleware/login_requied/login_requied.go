package login_requied

import (
	"08_ihome/apiv1.0/pkg/e"
	"08_ihome/apiv1.0/pkg/zaplog"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

var store = sessions.NewCookieStore([]byte("something-very-secret"))

func LoginRequied() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := store.Get(c.Request,"login")
		if err != nil {
			log.Print(err)
		}
		mobile := session.Values["mobile"]
		if mobile == nil {
			code := e.ERROR_AUTH
			zaplog.Logger.Errorf("未登录: %s code: %d","未登录",code)
			c.JSON(http.StatusOK,gin.H{
				"code": code,
				"msg":"未登录！",
			})
			c.Abort()
		}
	}
}
