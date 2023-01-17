package question

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"zhihu/app/internal/service"
)

// NewAnswer 发布回答
func NewAnswer(c *gin.Context) {
	content := c.PostForm("content")

	//解析API参数得到问题id
	idString := c.Param("id")
	id, _ := strconv.Atoi(idString)

	//从token中取出当前用户id
	uid := c.MustGet("uid").(int)

	err := service.PublicAnswer(id, uid, content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
			"ok":   false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "public answer successfully",
		"ok":   "true",
	})
}

// NewCommentToAnswer 发表对回答的评论
func NewCommentToAnswer(c *gin.Context) {
	content := c.PostForm("content")

	//解析API参数得到问题id
	idString := c.Param("id")
	id, _ := strconv.Atoi(idString)

	//从token中取出当前用户id
	uid := c.MustGet("uid").(int)

	//解析URL参数，得到父id
	pidString := c.Query("parentId")
	pid, _ := strconv.Atoi(pidString)

	err := service.CommentTheAnswer(id, uid, pid, content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
			"ok":   false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "public comment successfully",
		"ok":   "true",
	})
}

// ReplyNew 进行回复
func ReplyNew(c *gin.Context) {
	content := c.PostForm("content")

	//解析API参数得到问题id
	idString := c.Param("id")
	id, _ := strconv.Atoi(idString)

	//从token中取出当前用户id
	uid := c.MustGet("uid").(int)

	//解析URL参数，得到父id和被回复用户id
	pidString := c.Query("parentId")
	pid, _ := strconv.Atoi(pidString)
	toUidString := c.Query("toUid")
	toUid, _ := strconv.Atoi(toUidString)

	err := service.ReplyToComment(id, uid, pid, toUid, content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
			"ok":   false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "reply successfully",
		"ok":   "true",
	})
}
