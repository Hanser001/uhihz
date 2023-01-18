package router

import (
	"github.com/gin-gonic/gin"
	"zhihu/app/api/article"
	"zhihu/app/api/question"
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
		users.PUT("/edit/username", middleware.JWTAuthMiddleware(), user.Usernames)
		users.PUT("/edit/personalSignature", middleware.JWTAuthMiddleware(), user.PersonalSignature)
	}

	articles := v1.Group("articles")
	{
		articles.POST("/new", middleware.JWTAuthMiddleware(), article.NewArticle)
		articles.PUT("/new/:id", middleware.JWTAuthMiddleware(), article.UpdateArticle)
		articles.GET("/content/:id", article.ReadArticle)
		articles.POST("/:id/comment", middleware.JWTAuthMiddleware(), article.NewComment)
		articles.POST("/:id/comment/comment", middleware.JWTAuthMiddleware(), article.NewCommentToParentComment)
		articles.POST("/:id/reply", middleware.JWTAuthMiddleware(), article.NewReply)
		articles.POST("/:id/like", middleware.JWTAuthMiddleware(), article.NewComment)
	}

	questions := v1.Group("questions")
	{
		questions.POST("/new", middleware.JWTAuthMiddleware(), question.NewQuestion)
		questions.PUT("/new/:id", middleware.JWTAuthMiddleware(), question.UpdateQuestion)
		questions.GET("/content/:id", question.ReadQuestion)
		questions.POST("/:id/answer", middleware.JWTAuthMiddleware(), question.NewAnswer)
		questions.POST("/:id/comment", middleware.JWTAuthMiddleware(), question.NewCommentToAnswer)
		questions.POST("/:id/reply", middleware.JWTAuthMiddleware(), question.ReplyNew)
	}

	return r
}
