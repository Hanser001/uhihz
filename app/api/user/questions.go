package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"zhihu/app/internal/service"
)

// NewQuestion 发布新问题
func NewQuestion(c *gin.Context) {
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

	//从token中取得作者id
	uid := c.MustGet("uid").(int)

	service.PublicQuestion(uid, title, content)
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "public question successfully",
		"ok":   true,
	})
}

// UpdateQuestion 更新问题
func UpdateQuestion(c *gin.Context) {
	newContent := c.PostForm("content")

	//这个id应该是问题id，但如何验证当前用户身份？
	idString := c.Param("id")
	id, _ := strconv.Atoi(idString)

	service.UpdateQuestion(newContent, id)

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "update question successfully",
		"ok":   true,
	})
}

// ReadQuestion 查看问题
func ReadQuestion(c *gin.Context) {
	//解析API参数得到问题id
	idString := c.Param("id")
	id, _ := strconv.Atoi(idString)

	question := service.SelectQuestion(id)

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "get the article",
		"data": question,
	})
}
