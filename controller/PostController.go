package controller

import (
	"log"
	"strconv"
	"study-gin-gorm/common"
	"study-gin-gorm/model"
	"study-gin-gorm/requestParam"
	"study-gin-gorm/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IPostController interface {
	RestController
	PageList(ctx *gin.Context)
}

type PostController struct {
	DB *gorm.DB
}

func NewPostController() IPostController {
	db := common.GetDB()
	db.AutoMigrate(model.Post{})

	return PostController{DB: db}
}

func (p PostController) PageList(ctx *gin.Context) {
	// 获取分页参数 如果未传入，默认显示 第1页，最多共20条
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// 分页
	var posts []model.Post
	p.DB.Preload("Category").Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&posts)

	// 前端渲染分页需要知道总数
	var total int64
	p.DB.Model(model.Post{}).Count(&total)

	response.Success(ctx, gin.H{"data": posts, "total": total}, "成功")

}

func (p PostController) Create(ctx *gin.Context) {
	// 绑定body 中的参数
	var requestPost requestParam.CreatePostRequest

	// 数据验证
	if err := ctx.ShouldBind(&requestPost); err != nil {
		log.Println(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// 获取登录用户 user
	user, _ := ctx.Get("user")

	// 创建文章 post
	post := model.Post{
		UserID:     user.(model.User).ID,
		CategoryID: requestPost.CategoryID,
		Title:      requestPost.Title,
		HeadImg:    requestPost.HeadImg,
		Content:    requestPost.Content,
	}

	if err := p.DB.Create(&post).Error; err != nil {
		panic(err)
		return
	}

	response.Success(ctx, nil, "创建成功")
}

func (p PostController) Update(ctx *gin.Context) {
	// 绑定body 中的参数
	var requestPost requestParam.CreatePostRequest

	// 数据验证
	if err := ctx.ShouldBind(&requestPost); err != nil {
		response.Fail(ctx, nil, "数据验证错误，分类名称必填")
		return
	}

	// 获取path 中的参数
	postID := ctx.Params.ByName("id")

	var post model.Post
	if p.DB.Where("id = ?", postID).First(&post).RowsAffected == 0 {
		response.Fail(ctx, nil, "文章不存在")
		return
	}

	// 判断当前用户是否为文章的作者
	// 获取登录用户 user
	user, _ := ctx.Get("user")
	userID := user.(model.User).ID
	if userID != post.UserID {
		response.Fail(ctx, nil, "文章不属于您，请勿非法操作")
		return
	}

	// 更新文章
	if err := p.DB.Model(&post).Updates(requestPost).Error; err != nil {
		response.Fail(ctx, nil, "更新失败")
		return
	}
	response.Success(ctx, gin.H{"post": post}, "更新成功")
}

func (p PostController) Show(ctx *gin.Context) {
	// 获取path 中的参数
	postID := ctx.Params.ByName("id")

	var post model.Post
	if p.DB.Preload("Category").Where("id = ?", postID).First(&post).RowsAffected == 0 {
		response.Fail(ctx, nil, "文章不存在")
		return
	}

	response.Success(ctx, gin.H{"post": post}, "成功")
}

func (p PostController) Delete(ctx *gin.Context) {
	// 获取path 中的参数
	postID := ctx.Params.ByName("id")

	var post model.Post
	if p.DB.Where("id = ?", postID).First(&post).RowsAffected == 0 {
		response.Fail(ctx, nil, "删除失败,文章不存在")
		return
	}

	// 判断当前用户是否为文章的作者
	// 获取登录用户 user
	user, _ := ctx.Get("user")
	userID := user.(model.User).ID
	if userID != post.UserID {
		response.Fail(ctx, nil, "文章不属于您，请勿非法操作")
		return
	}

	p.DB.Delete(&post)

	response.Success(ctx, gin.H{"post": post}, "删除成功")
}
