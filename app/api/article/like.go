package article

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"zhihu/app/internal/service"
)

// LikeActionToArticle  点赞(取消点赞)文章
func LikeActionToArticle(c *gin.Context) {
	//解析api参数，得到文章id
	idString := c.Param("id")
	id, _ := strconv.Atoi(idString)

	//从token中取出当前用户id
	uid := c.MustGet("uid").(int)

	//判断是否点赞
	flag := service.JudgeLikeToArticle(c, id, uid)

	if flag {
		//若点过赞，则拦截
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "have liked",
			"ok":   false,
		})
	} else if !flag {
		//若未点赞则点赞
		service.LikeToArticle(c, id, uid)
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "add the like!",
			"ok":   true,
		})
	}
}

// UnlikeActionToArticle 取消对文章点赞
func UnlikeActionToArticle(c *gin.Context) {
	//解析api参数，得到文章id
	idString := c.Param("id")
	id, _ := strconv.Atoi(idString)

	//从token中取出当前用户id
	uid := c.MustGet("uid").(int)

	//判断是否点赞
	flag := service.JudgeLikeToArticle(c, id, uid)

	if flag {
		//点过赞就取消点赞
		service.UnlikeToArticle(c, id, uid)
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

// LikeActionToComment  对评论点赞
func LikeActionToComment(c *gin.Context) {
	//解析url参数，得到评论（回复）id
	cidString := c.Query("cid")
	cid, _ := strconv.Atoi(cidString)

	//从token中取出当前用户id
	uid := c.MustGet("uid").(int)

	//判断是否点赞
	flag := service.JudgeLikeToComment(c, cid, uid)

	if flag {
		//若点过赞则拦截
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "have liked!",
			"ok":   false,
		})
		return
	} else if !flag {
		//若未点赞，则会点赞
		service.LikeArticleComment(c, cid, uid)
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "add the like!",
			"ok":   true,
		})
	}

}

// UnlikeActionToComment 取消对评论点赞
func UnlikeActionToComment(c *gin.Context) {
	//解析url参数，得到评论（回复）id
	cidString := c.Query("cid")
	cid, _ := strconv.Atoi(cidString)

	//从token中取出当前用户id
	uid := c.MustGet("uid").(int)

	//判断是否点赞
	flag := service.JudgeLikeToComment(c, cid, uid)

	if flag {
		//若点过赞则取消点赞
		service.UnlikeArticleComment(c, cid, uid)
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "cancel the like!",
			"ok":   true,
		})
	} else if !flag {
		//未点赞则拦截
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "not like!",
			"ok":   false,
		})
		return
	}
}
