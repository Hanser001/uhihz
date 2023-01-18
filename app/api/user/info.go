package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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

	c.JSON(http.StatusOK, data)
}

// GetUserAnswers 取得用户回答
func GetUserAnswers(c *gin.Context) {
	//从api参数解析出用户id
	uidString := c.Param("id")
	uid, _ := strconv.Atoi(uidString)

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
	//从api参数解析出用户id
	uidString := c.Param("id")
	uid, _ := strconv.Atoi(uidString)

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

// GetUserQuestions 取得用户提问
func GetUserQuestions(c *gin.Context) {
	//从api参数解析出用户id
	uidString := c.Param("id")
	uid, _ := strconv.Atoi(uidString)

	data, err := service.GetUserQuestions(uid)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
			"ok":   false,
		})
		return
	}

	c.JSON(http.StatusOK, data)
}
