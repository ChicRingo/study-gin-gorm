package router

import (
	"net/http"
	"study-gin-gorm/controller"
	"study-gin-gorm/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 跨域访问中间件
	r.Use(middleware.CORSMiddleware(), middleware.RecoveryMiddleware())

	// 注册路由：用户注册、登录、用户信息
	r.POST("/api/auth/register", controller.Register)
	r.POST("/api/auth/login", controller.Login)
	r.GET("/api/auth/info", middleware.JWTAuthMiddleware(), controller.Info)

	// 注册路由组：分类
	categoryRouter := r.Group("/categories")
	{
		categoryController := controller.NewCategoryController()
		categoryRouter.POST("", categoryController.Create)
		categoryRouter.PUT(":id", categoryController.Update)
		categoryRouter.GET(":id", categoryController.Show)
		categoryRouter.DELETE(":id", categoryController.Delete)
	}

	// 注册路由组：帖子
	postRouter := r.Group("/posts")
	postRouter.Use(middleware.JWTAuthMiddleware())
	{
		postController := controller.NewPostController()
		postRouter.POST("", postController.Create)
		postRouter.PUT(":id", postController.Update)
		postRouter.GET(":id", postController.Show)
		postRouter.DELETE(":id", postController.Delete)
		postRouter.POST("page/list", postController.PageList)
	}

	// 未定义路由组
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "404",
		})
	})

	return r
}
