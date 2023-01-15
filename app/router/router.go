package router

import (
	"github.com/gin-gonic/gin"
	"zhihu/app/api/user"
	g "zhihu/app/global"
	"zhihu/app/internal/middleware"
)

func InitRouters() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.ZapLogger(g.Logger), middleware.ZapRecovery(g.Logger, true), middleware.CORS())

	api := r.Group("/api")
	v1 := api.Group("/v1")
	users := v1.Group("users")

	users.POST("/register", user.Register)
	users.POST("login", user.Login)
	users.GET("/:id", user.PersonalInfo)
	users.PUT("/security/password", user.NewPassword)

	return r
}