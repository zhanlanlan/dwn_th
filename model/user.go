package model

import (
	"context"
	"dwn_th/db"
	"dwn_th/proto"
	"errors"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName      string `gorm:"column:user_name; size:25; unique; not null"`
	PassWord      string `gorm:"column:pass_word; size:255; not null"`
	LastLoginTime int64  `gorm:"column:last_login_time; not null; default:0"`
}

func (u *User) TableName() string {
	return "t_user"
}

func TUser(ctx context.Context) *gorm.DB {
	return db.DB.WithContext(ctx).Table("t_user")
}

func CreateUser(ctx context.Context, user *User) error {
	return TUser(ctx).Create(user).Error
}

func DeleteUser(ctx context.Context, userName string) error {
	return TUser(ctx).Where("user_name = ?", userName).Delete(&User{}).Error
}

func GetUserByUserName(ctx context.Context, userName string) (u User, err error) {
	err = TUser(ctx).Where("user_name = ?", userName).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = proto.UserNotFound
		return
	}

	return
}

func UpdateUserPassWord(ctx context.Context, userName string, passWord string) error {
	err := TUser(ctx).Where("user_name = ?", userName).Update("pass_word", passWord).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return proto.UserNotFound
	}

	return nil
}

func UpdateUserLoginTime(ctx context.Context, userId int64, t int64) error { return nil }
