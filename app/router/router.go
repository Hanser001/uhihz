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
		users.GET("/:id/answers", user.GetUserAnswers)
		users.GET("/:id/articles", user.GetUserArticles)
		users.GET("/:id/questions", user.GetUserQuestions)
		users.GET("/:id/collections", user.GetArticleCollection)
	}

	articles := v1.Group("articles")
	{
		articles.POST("/new", middleware.JWTAuthMiddleware(), article.NewArticle)
		articles.PUT("/new/:id", middleware.JWTAuthMiddleware(), article.UpdateArticle)
		articles.DELETE("/:id/deleted", middleware.JWTAuthMiddleware(), article.DeleteArticle)
		articles.GET("/content/:id", article.ReadArticle)
		articles.POST("/:id/like", middleware.JWTAuthMiddleware(), article.LikeActionToArticle)
		articles.DELETE("/:id/like", middleware.JWTAuthMiddleware(), article.UnlikeActionToArticle)
		articles.POST("/:id/collection", middleware.JWTAuthMiddleware(), article.CollectActionToArticle)
		articles.DELETE("/:id/collection", middleware.JWTAuthMiddleware(), article.CancelActionToArticle)

		articles.POST("/:id/comment", middleware.JWTAuthMiddleware(), article.NewComment)
		articles.POST("/:id/comment/review", middleware.JWTAuthMiddleware(), article.NewCommentToParentComment)
		articles.POST("/:id/reply", middleware.JWTAuthMiddleware(), article.NewReply)
		articles.GET("/comment", article.GetComment)
		articles.DELETE("/:id/comment/deleted", middleware.JWTAuthMiddleware(), article.DeleteComment)
		articles.POST("/:id/comment/like", middleware.JWTAuthMiddleware(), article.LikeActionToComment)
		articles.DELETE("/:id/comment/like", middleware.JWTAuthMiddleware(), article.UnlikeActionToComment)
	}

	questions := v1.Group("questions")
	{
		questions.POST("/new", middleware.JWTAuthMiddleware(), question.NewQuestion)
		questions.PUT("/new/:id", middleware.JWTAuthMiddleware(), question.UpdateQuestion)
		questions.GET("/content/:id", question.ReadQuestion)
		questions.DELETE("/:id/deleted", middleware.JWTAuthMiddleware(), question.DeleteQuestion)
		questions.POST("/:id/like", middleware.JWTAuthMiddleware(), question.LikeActionToQuestion)
		questions.DELETE("/:id/like", middleware.JWTAuthMiddleware(), question.UnlikeActionToQuestion)

		questions.POST("/:id/answer", middleware.JWTAuthMiddleware(), question.NewAnswer)
		questions.POST("/:id/comment", middleware.JWTAuthMiddleware(), question.NewCommentToAnswer)
		questions.POST("/:id/reply", middleware.JWTAuthMiddleware(), question.ReplyNew)
		questions.DELETE("/comment/deleted", middleware.JWTAuthMiddleware(), question.DeleteComment)
		questions.GET("/comment", question.GetComment)
		questions.POST("/comment/like", middleware.JWTAuthMiddleware(), question.LikeActionToComment)
		questions.DELETE("/comment/like", middleware.JWTAuthMiddleware(), question.UnlikeActionToComment)
		questions.PUT("/comment/new", middleware.JWTAuthMiddleware(), question.UpdateAnswer)
	}

	return r
}
