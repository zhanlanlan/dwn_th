package model

import (
	"dwn_th/utils"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 全局数据库链接
var DB *gorm.DB

func InitDB() {
	if DB == nil {

		dsn := utils.MustGetENV("dsn")

		var err error
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("初始化数据库链接失败: %s", err.Error())
		}

		if utils.OptGetEnvBool("DEBUG") {
			DB = DB.Debug()
		}
	}
}

func Migerate() {
	DB.AutoMigrate(&User{}, &File{}, &UserFile{})
}
