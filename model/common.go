package model

import "dwn_th/db"

func Migerate() {
	db.DB.AutoMigrate(&User{}, &File{}, &UserFile{})
}
