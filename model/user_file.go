package model

import (
	"context"
	"dwn_th/db"
	"dwn_th/proto"
	"errors"

	"gorm.io/gorm"
)

type UserFile struct {
	gorm.Model
	UserId   int64  `gorm:"column:user_id; not null"`
	FileId   int64  `gorm:"column:file_id; not null"`
	Pwd      string `gorm:"column:pwd; size:1024; not null"`
	FileName string `gorm:"column:file_name; size:255; not null; unique"`
	FileType int    `gorm:"column:file_type; not null"`
	Ext      string `gorm:"column:ext; size:50; default:"`
}

func (f *UserFile) TableName() string {
	return "t_user_file"
}

func TUserFile(ctx context.Context) *gorm.DB {
	return db.DB.WithContext(ctx).Table("t_user_file")
}

func GetFileList(ctx context.Context, pwd string, userID int64) (userfiles []UserFile, err error) {
	userfiles = make([]UserFile, 0)
	err = TUserFile(ctx).Where("pwd = ? AND user_id = ?", pwd, userID).Find(&userfiles).Error
	return
}

func GetFile(ctx context.Context, pwd string, name string) (f UserFile, err error) {
	err = TUserFile(ctx).Where("pwd = ? AND file_name = ? ", pwd, name).First(&f).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = proto.FileNotFound
		return
	}
	return

}

func CreateUserFile(ctx context.Context, uf *UserFile) error {
	return TUserFile(ctx).Create(uf).Error
}

func GetFileByOwner(ctx context.Context, owner int64, pwd string, fileName string) (f UserFile, err error) {
	err = TUserFile(ctx).Where("user_id = ? AND pwd = ? AND file_name = ?", owner, pwd, fileName).First(&f).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = proto.FileNotFound
		return
	}

	return
}

func GetDir(ctx context.Context, pwd, name string) (f UserFile, err error) {
	err = TUserFile(ctx).Where("pwd = ? AND file_name = ? AND file_type = ?", pwd, name, TypeDir).First(&f).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = proto.DirNotExist
		return
	}
	return
}

func DeleteFile(ctx context.Context, owner int64, pwd string, fileName string) error {
	return TUserFile(ctx).Where("user_id = ? AND pwd = ? AND file_name = ?", owner, pwd, fileName).Delete(&UserFile{}).Error
}
