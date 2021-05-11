package db

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 全局数据库链接
var DB *gorm.DB

// TODO: 数据库用户名 密码 配置信息 应该写入配置文件
const dsn = "root:123456@tcp(127.0.0.1:3306)/start?charset=utf8mb4&parseTime=True&loc=Local"

func InitDB() {
	if DB == nil {
		var err error
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("初始化数据库链接失败: %s", err.Error())
		}

		// FIXME: 正式情况下应该删除
		DB = DB.Debug()
	}
}
