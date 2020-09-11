package controller

import (
	"fmt"
	"log"
	"strconv"
	"study-gin-gorm/model"
	"study-gin-gorm/repository"
	"study-gin-gorm/requestParam"
	"study-gin-gorm/response"

	"github.com/gin-gonic/gin"
)

type ICategoryController interface {
	RestController
}

type CategoryController struct {
	Repository repository.CategoryRepository
}

func NewCategoryController() ICategoryController {
	categoryRepository := repository.NewCategoryRepository()
	err := categoryRepository.DB.AutoMigrate(model.Category{})
	if err != nil {
		fmt.Println(err)
	}
	return CategoryController{Repository: categoryRepository}
}

func (c CategoryController) Create(ctx *gin.Context) {
	var requestCategory requestParam.CreateCategoryRequest
	if err := ctx.ShouldBind(&requestCategory); err != nil {
		log.Println(err.Error())
		response.Fail(ctx, nil, "数据验证错误，分类名称必填"+err.Error())
		return
	}

	category, err := c.Repository.Create(requestCategory.Name)
	if err != nil {
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	response.Success(ctx, gin.H{"category": category}, "创建成功")
}

func (c CategoryController) Update(ctx *gin.Context) {
	// 绑定body 中的参数
	var requestCategory requestParam.CreateCategoryRequest
	if err := ctx.ShouldBind(&requestCategory); err != nil {
		response.Fail(ctx, nil, "数据验证错误，分类名称必填")
		return
	}

	// 获取path 中的参数
	categoryID, _ := strconv.Atoi(ctx.Params.ByName("id"))

	updateCategory, err := c.Repository.SelectByID(categoryID)
	if err != nil {
		response.Fail(ctx, nil, "分类不存在")
		return
	}

	// 更新分类
	// map
	// struct
	// name value
	category, err := c.Repository.Update(*updateCategory, requestCategory.Name)
	if err != nil {
		panic(err)
	}

	response.Success(ctx, gin.H{"category": category}, "修改成功")
}

func (c CategoryController) Show(ctx *gin.Context) {
	// 获取path 中的参数
	categoryID, _ := strconv.Atoi(ctx.Params.ByName("id"))

	category, err := c.Repository.SelectByID(categoryID)
	if err != nil {
		response.Fail(ctx, nil, "分类不存在")
		return
	}

	response.Success(ctx, gin.H{"category": category}, "")
}

func (c CategoryController) Delete(ctx *gin.Context) {
	// 获取path 中的参数
	categoryID, _ := strconv.Atoi(ctx.Params.ByName("id"))

	if err := c.Repository.DeleteByID(categoryID); err != nil {
		response.Fail(ctx, nil, "删除失败，请重试")
		return
	}

	response.Success(ctx, nil, "删除成功")
}
