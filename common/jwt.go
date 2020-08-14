package common

import (
	"study-gin-gorm/model"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("a_secret_crect")

type Clamis struct {
	UserID uint
	jwt.StandardClaims
}

func ReleaseToken(user model.User) (string, error) {
	//有效时间
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Clamis{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "demo.studygin",
			Subject:   "user token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenString string) (*jwt.Token, *Clamis, error) {
	clamis := &Clamis{}
	token, err := jwt.ParseWithClaims(tokenString, clamis, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	return token, clamis, err
}
