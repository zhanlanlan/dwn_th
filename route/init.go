package route

import (
	"dwn_th/db"
	"dwn_th/model"
	"dwn_th/storage"
)

func InitClients() {
	db.InitDB()
	model.Migerate()

	storage.InitOssClient()
}
