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
	{
		users.POST("/register", user.Register)
		users.POST("login", user.Login)
		users.GET("/:id", user.PersonalInfo)
		users.PUT("/security/password", user.NewPassword)
		users.PUT("/info/username", middleware.JWTAuthMiddleware(), user.Usernames)
		users.PUT("/info/personalSignature", middleware.JWTAuthMiddleware(), user.PersonalSignature)
	}

	articles := v1.Group("articles")
	{
		articles.POST("/new", middleware.JWTAuthMiddleware(), user.NewArticle)
		articles.PUT("/new/:id", middleware.JWTAuthMiddleware(), user.UpdateArticle)
		articles.GET("/content/:id", user.ReadArticle)
	}

	questions := v1.Group("questions")
	{
		questions.POST("/new", middleware.JWTAuthMiddleware(), user.NewQuestion)
		questions.PUT("/new/:id", middleware.JWTAuthMiddleware(), user.UpdateQuestion)
		questions.GET("/content/:id", user.ReadQuestion)
	}
	return r
}
