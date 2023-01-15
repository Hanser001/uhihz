package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

const TokenExpireDuration = time.Hour * 24

// CustomSecret 用于加盐的字符串
var CustomSecret = []byte("nachoNeko")

type CustomClaims struct {
	//额外保存用户id
	Uid int
	jwt.RegisteredClaims
}

// GenToken 生成JWT
func GenToken(uid int) (string, error) {
	// 创建一个我们自己的声明
	c := CustomClaims{
		uid, // 自定义字段
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)), // 过期时间
			Issuer:    "uhihz",                                                 // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(CustomSecret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*CustomClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return CustomSecret, nil
	})
	if err != nil {
		return nil, err
	}
	// 对token对象中的Claim进行类型断言
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
