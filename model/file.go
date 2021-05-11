package model

import (
	"context"
	"dwn_th/db"

	"gorm.io/gorm"
)

type FileType = int64

const (
	TypeFile FileType = 1
	TypeDir           = 2
)

type File struct {
	gorm.Model
	FileKey string `gorm:"column:file_key; size:255; not null; unique"`
	URI     string `gorm:"column:uri; size:50; default:"`
	Size    int64  `gorm:"column:size; not null; default:0"`
	Hash    string `gorm:"column:hash; size:64; not null"`
}

func (f *File) TableName() string {
	return "t_file"
}

func TFile(ctx context.Context) *gorm.DB {
	return db.DB.WithContext(ctx).Table("t_file")
}

func CreateFile(ctx context.Context, f *File) error {
	return TFile(ctx).Create(f).Error
}

func UpdateDownloadTimes(ctx context.Context, fileID int64) error {
	panic("unimplemented!")
}

func GetFileById(ctx context.Context, fileId int64) (file File, err error) {
	err = TFile(ctx).Where("id = ?", fileId).First(&file).Error
	return
}

func GetFileByHash(ctx context.Context, hash string) (file File, err error) {
	err = TFile(ctx).Where("Hash = ?", hash).First(&file).Error
	return
}

func GetFileByFileKey(ctx context.Context, filekey string) (file File, err error) {
	err = TFile(ctx).Where("file_key = ?", filekey).First(&file).Error
	return
}
