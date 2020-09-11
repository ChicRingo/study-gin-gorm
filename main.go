package main

import (
	"fmt"
	"study-gin-gorm/common"
	"study-gin-gorm/config"
	"study-gin-gorm/router"

	"github.com/spf13/viper"
)

func main() {
	// 1.加载配置文件
	if err := config.Init(); err != nil {
		fmt.Printf("init settings failed, err:%v\n", err)
		return
	}

	// 3.初始化MySQL
	if err := common.Init(config.Conf.MySQLConfig); err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	defer common.Close()

	//// 注册gin框架内置的校验器翻译
	//if err := controller.InitTrans("zh"); err != nil {
	//	fmt.Printf("init validator failed, err:%v\n", err)
	//}

	// 注册路由
	r := router.SetupRouter()
	port := viper.GetString("server.port")
	if port != "" {
		panic(r.Run(":" + port))
	}
	panic(r.Run()) // 监听并在 0.0.0.0:8080 上启动服务
}
