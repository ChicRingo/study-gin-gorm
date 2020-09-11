package jwt

import (
	"errors"
	"study-gin-gorm/model"
	"time"

	"github.com/spf13/viper"

	"github.com/dgrijalva/jwt-go"
)

// secret签名
var jwtSecretKey = []byte("a_jwt_secret_key")

/*
自定义声明结构体并内嵌jwt.StandardClaims
jwt包自带的 jwt.StandardClaims 只包含了官方字段
我们这里需要额外记录 user_id 字段，所以要自定义结构体
如果想要保存更多信息，都可以添加到这个结构体中
*/
type Claims struct {
	UserID uint // 自定义字段
	jwt.StandardClaims
}

// 生成token
func ReleaseToken(user model.User) (tokenStr string, err error) {
	// 创建一个自定义的声明
	claims := &Claims{
		user.ID,
		jwt.StandardClaims{
			// 过期时间
			ExpiresAt: time.Now().Add(viper.GetDuration("auth.jwt_expire") * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
			// 签发人
			Issuer:  "study_gin_demo",
			Subject: "user token",
		},
	}

	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用指定的secret签名并获得完整的编码后的字符串token
	tokenStr, err = token.SignedString(jwtSecretKey)
	return
}

// 解析token
func ParseToken(tokenString string) (claims *Claims, err error) {
	claims = new(Claims)
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})

	// 校验token
	if err != nil || !token.Valid {
		err = errors.New("invalid token")
		return nil, err
	}

	return claims, nil
}
