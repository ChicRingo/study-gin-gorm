package controller

import (
	"log"
	"net/http"
	"study-gin-gorm/dao/mysql"
	"study-gin-gorm/dto"
	"study-gin-gorm/model"
	"study-gin-gorm/pkg/jwt"
	"study-gin-gorm/pkg/util"
	"study-gin-gorm/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

//注册
func Register(ctx *gin.Context) {
	db := mysql.GetDB()
	// 使用map 获取请求的参数
	//var requestMap =  make(map[string]string)
	//json.NewDecoder(ctx.Request.Body).Decode(&requestMap)

	// 使用结构体
	//var requestUser = model.User{}
	//json.NewDecoder(ctx.Request.Body).Decode(&requestUser)

	//gin框架Bind函数
	var requestUser = model.User{}
	if err := ctx.ShouldBindJSON(&requestUser); err != nil {
		// 请求参数有误，直接返回响应

		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			response.Fail(ctx, nil, "参数校验错误，请重试")
			return
		}
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, errs.Error())
		return
	}

	//获取参数
	//name := ctx.PostForm("name")
	//telephone := ctx.PostForm("telephone")
	//password := ctx.PostForm("password")
	name := requestUser.Name
	telephone := requestUser.Telephone
	password := requestUser.Password

	////数据验证
	//if len(telephone) != 11 {
	//	response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位")
	//	return
	//}
	if len(password) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}
	if len(name) == 0 {
		name = util.RandomString(10)
	}

	log.Println(name, telephone, password)

	//判断手机号是否存在
	if isTelephoneExist(db, telephone) {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户已经存在")
		return
	}

	//创建用户
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "加密错误")
		return
	}
	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hasedPassword),
	}
	db.Create(&newUser)

	//发放token
	token, err := jwt.ReleaseToken(newUser)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "系统异常")
		log.Printf("token生成异常，错误信息：%v", err)
		return
	}

	//返回结果
	response.Success(ctx, gin.H{"token": token}, "注册成功")
}

//登录
func Login(ctx *gin.Context) {
	db := mysql.GetDB()
	//gin框架Bind函数
	var requestUser = model.User{}
	ctx.ShouldBindJSON(&requestUser)

	//获取参数
	telephone := requestUser.Telephone
	password := requestUser.Password

	//数据验证
	if len(telephone) != 11 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位")
		return
	}
	if len(password) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}

	//判断手机号是否存在
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户不存在")
		//ctx.JSON(http.StatusUnprocessableEntity, gin.H{
		//	"code": 422,
		//	"msg":  "用户不存在",
		//})
		return
	}

	//判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		response.Response(ctx, http.StatusBadRequest, 400, nil, "密码错误")
		//ctx.JSON(http.StatusBadRequest, gin.H{
		//	"code": 400,
		//	"msg":  "密码错误",
		//})
		return
	}

	//发放token
	token, err := jwt.ReleaseToken(user)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "系统异常")
		//ctx.JSON(http.StatusInternalServerError, gin.H{
		//	"code": 500,
		//	"msg":  "系统异常",
		//})
		log.Printf("token生成异常，错误信息：%v", err)
		return
	}
	//返回结果
	response.Success(ctx, gin.H{"token": token}, "登录成功")
}

func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	response.Success(ctx, gin.H{"user": dto.ToUserDto(user.(model.User))}, "")
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"user": dto.ToUserDto(user.(model.User)),
		},
	})
}

//判断手机号是否存在
func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}

	return false
}
