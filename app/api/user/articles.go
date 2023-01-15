package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"zhihu/app/internal/service"
)

// NewArticle 发布新文章
func NewArticle(c *gin.Context) {
	title := c.PostForm("title")
	content := c.PostForm("content")

	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "title cannot be null",
			"ok":   false,
		})
		return
	}

	if content == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "content cannot be null",
			"ok":   false,
		})
		return
	}

	//从token中取得作者id
	uid := c.MustGet("uid").(int)

	service.PublishArticle(uid, title, content)
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "public article successfully",
		"ok":   true,
	})
}

// UpdateArticle 更新文章
func UpdateArticle(c *gin.Context) {
	newContent := c.PostForm("content")

	if newContent == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "content can not be null",
			"ok":   false,
		})
	}

	//这个id应该是文章id，但如何验证当前用户身份？
	idString := c.Param("id")
	id, _ := strconv.Atoi(idString)

	service.UpdateArticle(newContent, id)
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "update article successfully",
		"ok":   true,
	})
}

// ReadArticle 阅读文章
func ReadArticle(c *gin.Context) {
	//解析API参数，得到阅读文章id
	idString := c.Param("id")
	id, _ := strconv.Atoi(idString)

	article := service.SelectArticle(id)

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "get the article",
		"data": article,
	})
}
