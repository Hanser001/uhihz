package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zhihu/app/internal/service"
)

// GetArticleCollection 获取用户收藏文章
func GetArticleCollection(c *gin.Context) {
	//从token取出uid
	uid := c.MustGet("uid").(int)

	data, err := service.GetUserArticleCollection(uid)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
			"ok":   false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

// GetUserAnswers 取得用户回答
func GetUserAnswers(c *gin.Context) {
	//从token取出uid
	uid := c.MustGet("uid").(int)

	data, err := service.GetUserAnswer(uid)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
			"ok":   false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

// GetUserArticles 取得用户发布文章
func GetUserArticles(c *gin.Context) {
	//从token取出uid
	uid := c.MustGet("uid").(int)

	data, err := service.GetUserArticles(uid)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
			"ok":   false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}
