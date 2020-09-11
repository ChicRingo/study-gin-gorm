package jwt

import (
	"study-gin-gorm/model"
	"time"

	"github.com/spf13/viper"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecretKey = []byte("a_jwt_secret_key")

type Claims struct {
	UserID uint
	jwt.StandardClaims
}

func ReleaseToken(user model.User) (string, error) {
	// 创建一个自定义的声明数据
	claims := &Claims{
		user.ID, // 自定义字段
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
	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// 解析token
func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := new(Claims)
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})
	return token, claims, err
}
