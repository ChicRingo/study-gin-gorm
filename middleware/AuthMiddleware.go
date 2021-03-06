package middleware

import (
	"net/http"
	"strings"
	"study-gin-gorm/dao/mysql"
	"study-gin-gorm/model"
	"study-gin-gorm/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//获取authorization header
		tokenString := ctx.GetHeader("Authorization")

		//验证token格式
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": "401",
				"msg":  "权限不足",
			})
			ctx.Abort() //阻止调用中间件后续的函数
			return
		}

		tokenString = tokenString[7:]

		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
			ctx.Abort() //阻止调用中间件后续的函数
			return
		}

		//验证通过后获取claim中的userID
		userID := claims.UserID
		db := mysql.GetDB()
		var user model.User
		db.First(&user, userID)

		//用户
		if user.ID == 0 {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
			ctx.Abort() //阻止调用中间件后续的函数
			return
		}

		ctx.Set("user", user) //用户存在 将user的信息写入上下文

		ctx.Next() //验证通过，执行之后函数
	}
}
