package question

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"zhihu/app/internal/service"
)

// LikeActionToQuestion 点赞问题
func LikeActionToQuestion(c *gin.Context) {
	//解析api参数，得到问题id
	idString := c.Param("id")
	id, _ := strconv.Atoi(idString)

	//从token中取出用户id
	uid := c.MustGet("uid").(int)

	//判断是否点赞
	flag := service.JudgeLikeToQuestion(c, id, uid)

	if flag {
		//若点过赞,则拦截
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "have liked",
			"ok":   false,
		})
		return
	} else if !flag {
		//若未点赞,点击点赞则点赞
		service.LikeToQuestion(c, id, uid)
		//更新问题被点赞数量
		service.UpdateQuestionLikeNum(c, id, 1)
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "add the like!",
			"ok":   true,
		})
	}
}

// UnlikeActionToQuestion 取消对问题点赞
func UnlikeActionToQuestion(c *gin.Context) {
	//解析api参数，得到问题id
	idString := c.Param("id")
	id, _ := strconv.Atoi(idString)

	//从token中取出用户id
	uid := c.MustGet("uid").(int)

	//判断是否点赞
	flag := service.JudgeLikeToQuestion(c, id, uid)

	if flag {
		//点过赞就取消点赞
		service.UnlikeToQuestion(c, id, uid)
		//更新问题被点赞数量
		service.UpdateQuestionLikeNum(c, id, -1)
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "canceled",
			"ok":   true,
		})
	} else if !flag {
		//没点过赞就拦截
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "not like",
			"ok":   false,
		})
		return
	}
}

// LikeActionToComment 对评论点赞
func LikeActionToComment(c *gin.Context) {
	//从token中取出用户id
	uid := c.MustGet("uid").(int)
	//解析url参数，得到评论（回复）id
	cidString := c.Query("cid")
	cid, _ := strconv.Atoi(cidString)
	//判断是否点赞
	flag := service.JudgeLikeToAnswer(c, cid, uid)

	if flag {
		//点过赞则拦截
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "have liked",
			"ok":   false,
		})
		return
	} else if !flag {
		//没点赞就点赞
		service.LikeToAnswer(c, cid, uid)
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "add the like!",
			"ok":   true,
		})
	}
}

// UnlikeActionToComment 取消对评论点赞
func UnlikeActionToComment(c *gin.Context) {
	//从token中取出用户id
	uid := c.MustGet("uid").(int)
	//解析url参数，得到评论（回复）id
	cidString := c.Query("cid")
	cid, _ := strconv.Atoi(cidString)
	//判断是否点赞
	flag := service.JudgeLikeToAnswer(c, cid, uid)

	if flag {
		//点过赞则取消点赞
		service.UnlikeToAnswer(c, cid, uid)
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "add the like!",
			"ok":   true,
		})
	} else if !flag {
		//没点过赞则拦截
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "not like",
			"ok":   false,
		})
	}
}
