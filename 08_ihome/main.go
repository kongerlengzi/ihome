package main

import (
	_ "08_ihome/apiv1.0/model"
	"08_ihome/apiv1.0/pkg/zaplog"
	"08_ihome/apiv1.0/router"
	"github.com/gin-contrib/pprof"
)

func main()  {
	r := router.InitRouter()
	pprof.Register(r)
	zaplog.Init()
	//s := &http.Server{
	//	Addr: fmt.Sprintf(":%d", setting.HTTPPort),
	//	Handler: router,
	//	ReadTimeout: setting.ReadTimeout,
	//	WriteTimeout: setting.WriteTimeout,
	//	MaxHeaderBytes: 1 << 20,
	//}
	//
	//s.ListenAndServe()
	_ = r.Run()
}
