package boot

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	g "zhihu/app/global"
	"zhihu/app/router"
)

func ServerSet() {
	config := g.Config.Server

	gin.SetMode(config.Mode)
	routers := router.InitRouters()

	server := &http.Server{
		Addr:              config.Addr(),
		ReadTimeout:       config.GetReadTimeout(),
		WriteTimeout:      config.GetReadTimeout(),
		Handler:           routers,
		ReadHeaderTimeout: 0,
		IdleTimeout:       0,
		MaxHeaderBytes:    1 << 20, // 16 MB
	}

	g.Logger.Info("initialize server successfully!", zap.String("port", config.Addr()))
	g.Logger.Error(server.ListenAndServe().Error())
}
