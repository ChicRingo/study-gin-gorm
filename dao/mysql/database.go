package mysql

import (
	"fmt"
	"net/url"
	"study-gin-gorm/config"
	"study-gin-gorm/model"

	"go.uber.org/zap"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init(cfg *config.MySQLConfig) (err error) {
	args := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.Charset,
		url.QueryEscape(cfg.Loc),
	)

	db, err = gorm.Open(mysql.Open(args), &gorm.Config{
		//GORM v2自动创建外键，手动关闭该功能
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		zap.L().Error("connect DB failed, err:%v\n", zap.Error(err))
		return
	}
	//自动迁移 创建users表
	_ = db.AutoMigrate(&model.User{})
	return
}

// 获取数据库连接的DB指针
func GetDB() *gorm.DB {
	return db
}

// 关闭数据库连接
func Close() {
	mysqlDB, _ := db.DB()
	_ = mysqlDB.Close()
}
