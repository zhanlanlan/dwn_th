package route

import (
	"dwn_th/model"
	"dwn_th/storage"
	"log"

	"github.com/joho/godotenv"
)

func InitClients() {
	er := godotenv.Load()
	if er != nil {
		log.Fatalf("加载oss参数失败:%s", er.Error())
	}

	model.InitDB()
	model.Migerate()

	storage.InitOssClient()
}
