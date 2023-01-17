package article

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"zhihu/app/internal/service"
)

// NewComment 发布评论（一级评论）
func NewComment(c *gin.Context) {
	content := c.PostForm("content")

	if content == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "content can not be null",
			"ok":   false,
		})
		return
	}
	//解析api参数，得到文章id
	idString := c.Param("id")
	id, _ := strconv.Atoi(idString)

	//从token中取出当前用户id
	uid := c.MustGet("uid").(int)

	err := service.PublicComment(id, uid, content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
			"ok":   false,
		})
		return
	}

	//评论一次就为文章评论数加一
	service.AddArticleCommentNum(c, id)

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "public comment successfully",
		"ok":   true,
	})
}

// NewCommentToParentComment 发布对一级评论的评论
func NewCommentToParentComment(c *gin.Context) {
	content := c.PostForm("content")

	if content == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "content can not be null",
			"ok":   false,
		})
	}

	//解析api参数，得到文章id
	idString := c.Param("id")
	id, _ := strconv.Atoi(idString)

	//从token中取出当前用户id
	uid := c.MustGet("uid").(int)

	//解析URL参数，得到父评论id
	pidString := c.Query("parentId")
	pid, _ := strconv.Atoi(pidString)

	err := service.CommentTheComment(id, uid, pid, content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusOK,
			"msg":  err.Error(),
			"ok":   false,
		})
		return
	}

	//评论一次就为文章评论数加一
	service.AddArticleCommentNum(c, id)

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "public comment successfully",
		"ok":   true,
	})
}

// NewReply 进行回复
func NewReply(c *gin.Context) {
	content := c.PostForm("content")

	if content == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "content can not be null",
			"ok":   false,
		})
	}

	//解析api参数，得到文章id
	idString := c.Param("id")
	id, _ := strconv.Atoi(idString)

	//从token中取出当前用户id
	uid := c.MustGet("uid").(int)

	//解析URL参数，得到父评论id和被评论用户id
	pidString := c.Query("parentId")
	pid, _ := strconv.Atoi(pidString)
	toUidString := c.Query("toUid")
	toUid, _ := strconv.Atoi(toUidString)

	err := service.ReplyTheComment(id, uid, pid, toUid, content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
			"ok":   false,
		})
		return
	}

	//评论一次就为文章评论数加一
	service.AddArticleCommentNum(c, id)

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "reply",
		"ok":   true,
	})
}
