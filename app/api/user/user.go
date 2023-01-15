package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"zhihu/app/internal/service"
)

// Usernames 更改用户名
func Usernames(c *gin.Context) {
	newUsername := c.PostForm("newUsername")

	if newUsername == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "username cannot be null",
			"ok":   false,
		})
		return
	}

	err := service.SelectUserExists(newUsername)
	if err != nil {
		if err.Error() == "user already exists" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  err.Error(),
				"ok":   false,
			})
		} else if err.Error() == "internal error" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"msg":  err.Error(),
				"ok":   false,
			})
		}
		return
	}

	//从token中得到uid,没有则跳转到登录页
	id := c.MustGet("uid").(int)

	service.ChangeUsername(newUsername, id)

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "set new username successfully",
		"ok":   true,
	})
}

// NewPassword 修改密码
func NewPassword(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	newPassword := c.PostForm("newPassword")

	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "username cannot be null",
			"ok":   false,
		})
		return
	}
	if password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "password cannot be null",
			"ok":   false,
		})
		return
	}

	//分别校验用户名与密码
	flag1 := service.CheckUsername(username)
	if !flag1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "invalid username",
			"ok":   false,
		})
		return
	}

	encryptedPassword := service.GetPassword(username)

	flag2 := service.CheckPassword(password, encryptedPassword)
	if !flag2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "invalid password",
			"ok":   false,
		})
		return
	}

	newEncryptedPassword := service.EncryptPasswordWithSalt(newPassword)

	service.ChangePassword(newEncryptedPassword, username)
	c.JSON(http.StatusOK, gin.H{

		"code": http.StatusOK,
		"msg":  "set new password successfully",
		"ok":   true,
	})
}

// PersonalSignature 设置个性签名
func PersonalSignature(c *gin.Context) {
	personalSignature := c.PostForm("personalSignature")

	//从token中取出uid，没有则跳转登录页
	id := c.MustGet("uid").(int)

	service.SetPersonalSign(personalSignature, id)

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "set personalSignature successfully",
		"ok":   true,
	})
}

// PersonalInfo  查看用户基本信息
func PersonalInfo(c *gin.Context) {

	//解析API参数得到id
	idString := c.Param("id")
	id, _ := strconv.Atoi(idString)

	user := service.GetUserInfo(id)

	c.JSON(http.StatusOK, user)
}
