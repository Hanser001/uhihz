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
	err := service.PublicQuestion(uid, title, content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "public question successfully",
		"ok":   true,
	})
}

// UpdateQuestion 更新问题
func UpdateQuestion(c *gin.Context) {
	newContent := c.PostForm("content")

	//解析api参数取得问题id
	idString := c.Param("id")
	id, _ := strconv.Atoi(idString)

	//从token中取出当前用户id,在查询文章作者id，进行比较
	uid := c.MustGet("uid").(int)
	questionerId := service.GetQuestionerId(id)

	if uid != questionerId {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "you are not the questioner",
			"ok":   false,
		})
		return
	}

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

	//每次发送请求就会增加一次问题点击量
	service.AddQuestionClick(c, id)

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "get the article",
		"data": question,
	})
}
