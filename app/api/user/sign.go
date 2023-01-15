package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zhihu/app/internal/service"
	"zhihu/utils/jwt"
)

// Register 注册
func Register(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

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

	err := service.SelectUserExists(username)
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

	encryptedPassword := service.EncryptPasswordWithSalt(password)
	if encryptedPassword == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
			"ok":   false,
		})
		return
	}

	service.AddUser(username, encryptedPassword)
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "register successfully",
		"ok":   true,
	})
}

// Login 登录
func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

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
			"msg":  "invalid username or password",
			"ok":   false,
		})
		return
	}

	encryptedPassword := service.GetPassword(username)

	flag2 := service.CheckPassword(password, encryptedPassword)
	if !flag2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "invalid password or password",
			"ok":   false,
		})
		return
	}

	//登录成功则颁发token
	uid := service.GetUid(username)

	tokenString, _ := jwt.GenToken(uid)

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "login successfully",
		"data": gin.H{"token": tokenString},
	})
}
