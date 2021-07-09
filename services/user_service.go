package services

import (
	"context"
	"crypto/sha256"
	"dwn_th/model"
	"dwn_th/proto"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
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

func Login(ctx context.Context, userName, passWord string) *proto.Base {
	passWord = UtilHashPass(passWord)

	u, err := model.GetUserByUserName(ctx, userName)
	if errors.Is(err, proto.UserNotFound) {
		log.Println("用户不存在")
		return proto.UserNotFound
	} else if err != nil {
		log.Printf("登录失败: %s", err.Error())
		return proto.InternalErr
	}

	if u.PassWord != passWord {
		return proto.WrongPassword
	}

	log.Printf("User: %s logined", u.UserName)
	token, err := GenerateToken(u)
	if err != nil {
		log.Printf("login failed: %s", err.Error())
		return proto.LiginFailed
	}

	return &proto.Base{Code: 200, Msg: "success", Data: gin.H{"token": token}}
}
