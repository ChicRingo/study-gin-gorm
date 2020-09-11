package jwt

import (
	"study-gin-gorm/model"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("a_secret_crect")

type Claims struct {
	UserID uint
	jwt.StandardClaims
}

func ReleaseToken(user model.User) (string, error) {
	//有效时间
	expirationTime := time.Now().Add(7 * 24 * time.Hour)

	// 创建一个自定义的声明数据
	claims := &Claims{
		UserID: user.ID, // 自定义字段
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(), // 过期时间
			IssuedAt:  time.Now().Unix(),
			Issuer:    "study_gin_demo", // 签发人
			Subject:   "user token",
		},
	}

	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	return token, claims, err
}
