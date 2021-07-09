package services

import (
	"context"
	"crypto/sha256"
	"dwn_th/model"
	"dwn_th/proto"
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/golang/glog"
)

func UtilHashPass(pass string) string {
	sum := sha256.Sum256([]byte(pass))

	return base64.StdEncoding.EncodeToString(sum[:])
}

func CreateUser(ctx context.Context, req proto.CreateUserREQ) (err error) {
	req.PassWord = UtilHashPass(req.PassWord)

	_, err = model.GetUserByUserName(ctx, req.UserName)
	if !errors.Is(err, proto.UserNotFound) {
		return proto.UserAlreadyExist
	}

	u := model.User{
		UserName: req.UserName,
		PassWord: req.PassWord,
	}

	err = model.CreateUser(ctx, &u)
	if err != nil {
		return proto.InternalErr
	}

	return
}

func UpdatePsssword(ctx context.Context, userName, newPass string) (err error) {
	newPass = UtilHashPass(newPass)

	err = model.UpdateUserPassWord(ctx, userName, newPass)
	if errors.Is(err, proto.UserNotFound) {
		return proto.UserNotFound
	}

	return
}

type Claim struct {
	User model.User
	jwt.StandardClaims
}

func GenerateToken(user model.User) (token string, err error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, Claim{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 6).Unix(),
		},
	})

	return t.SignedString([]byte("fuck"))
}

func ParseToken(token string) (u model.User, err error) {
	data, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &u)
	return
}

func Login(ctx context.Context, userName, passWord string) (ret proto.LoginRES, err error) {
	passWord = UtilHashPass(passWord)

	u, err := model.GetUserByUserName(ctx, userName)
	if errors.Is(err, proto.UserNotFound) {
		glog.Errorf("用户 [%s] 不存在", userName)
		err = proto.UserNotFound
		return
	} else if err != nil {
		glog.Errorf("model.GetUserByUserName err: %s", err.Error())
		return
	}

	if u.PassWord != passWord {
		glog.Errorf("用户 [%s] 密码 [%s] 错误", userName, passWord)
		err = proto.WrongPassword
		return
	}

	glog.Infof("用户 [%s] 登录成功", userName)
	token, err := GenerateToken(u)
	if err != nil {
		glog.Errorf("GenerateToken err: %s", err.Error())
		return
	}

	ret.Token = token

	return
}
