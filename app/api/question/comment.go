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
	if content == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "content can not be null",
			"ok":   false,
		})
		return
	}

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

	//有人发布回答就增加问题的回答数
	service.UpdateQuestionAnswerNum(c, id, 1)

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "public answer successfully",
		"ok":   "true",
	})
}

// UpdateAnswer 更新回复
func UpdateAnswer(c *gin.Context) {
	content := c.PostForm("content")
	if content == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "content can not be null",
			"ok":   false,
		})
		return
	}

	//取得回复id
	cidString := c.Query("cid")
	cid, _ := strconv.Atoi(cidString)

	//将用户id和答者id取出，进行校验
	uid := c.MustGet("uid").(int)
	err, aid := service.GetAnswererId(cid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
			"ok":   false,
		})
		return
	}

	if uid != aid {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "no auth",
			"ok":   false,
		})
		return
	}

	err = service.UpdateAnswer(content, cid)
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
		"msg":  "update answer successfully",
		"ok":   "true",
	})
}

// NewCommentToAnswer 发表对回答的评论
func NewCommentToAnswer(c *gin.Context) {
	content := c.PostForm("content")
	if content == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "content can not be null",
			"ok":   false,
		})
	}

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

	if content == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "content can not be null",
			"ok":   false,
		})
	}

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

// DeleteComment 删除回答（回复）
func DeleteComment(c *gin.Context) {
	//解析url参数得到回复id
	cidString := c.Query("cid")
	cid, _ := strconv.Atoi(cidString)

	//取得答者id与当前用户id，进行校验
	err, rid := service.GetAnswererId(cid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
			"ok":   false,
		})
		return
	}

	uid := c.MustGet("uid").(int)

	if uid != rid {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "no auth",
			"ok":   false,
		})
		return
	}

	err = service.DeleteReview(cid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
			"ok":   false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "deleted",
		"ok":   true,
	})
}

// GetComment 查看回复
func GetComment(c *gin.Context) {
	//解析url参数得到回复id
	cidString := c.Query("cid")
	cid, _ := strconv.Atoi(cidString)

	err, reviewInfo := service.ReadReview(cid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
			"ok":   false,
		})
		return
	}

	c.JSON(http.StatusOK, reviewInfo)
}
