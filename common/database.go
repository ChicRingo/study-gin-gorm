package common

import (
	"fmt"
	"net/url"
	"study-gin-gorm/config"
	"study-gin-gorm/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() *gorm.DB {
	args := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=%s",
		config.Conf.MySQLConfig.Username,
		config.Conf.MySQLConfig.Password,
		config.Conf.MySQLConfig.Host,
		config.Conf.MySQLConfig.Port,
		config.Conf.MySQLConfig.Database,
		config.Conf.MySQLConfig.Charset,
		url.QueryEscape(config.Conf.MySQLConfig.Loc))

	//GORM v2自动创建外键，手动关闭该功能
	db, err := gorm.Open(mysql.Open(args), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic("连接数据库失败，错误信息：" + err.Error())
	}

	//自动迁移 创建users表
	_ = db.AutoMigrate(&model.User{})
	DB = db
	return db
}

func GetDB() *gorm.DB {
	return DB
}

func Close() {
	mysqlDB, _ := DB.DB()
	mysqlDB.Close()
}
