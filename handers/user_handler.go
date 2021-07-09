package handers

import (
	"log"

	"dwn_th/proto"
	"dwn_th/services"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

func CreateUser(c *gin.Context) {
	var req proto.CreateUserREQ
	if err := c.BindJSON(&req); err != nil {
		glog.Errorf("反序列化参数失败 err: %s", err.Error())
		proto.Err(c, err)
		return
	}

	if req.UserName == "" || len([]rune(req.UserName)) > 25 {
		proto.Wrap(c, nil, proto.BadUserName)
		return
	}

	if req.PassWord == "" {
		proto.Wrap(c, nil, proto.BadPassword)
		return
	}

	err := services.CreateUser(c, req)

	proto.Wrap(c, nil, err)
}

func UpdatePsssword(c *gin.Context) {
	var req proto.UpdatePssswordREQ
	if err := c.BindJSON(&req); err != nil {
		glog.Errorf("反序列化参数失败 err: %s", err.Error())
		proto.Err(c, proto.BadRquest)
		return
	}

	claim := Who(c)

	if req.NewPassWord == "" {
		//
		return
	}

	err := services.UpdatePsssword(c, claim.User.UserName, req.NewPassWord)
	proto.Wrap(c, nil, err)
}

func Login(c *gin.Context) {
	var req proto.LoginREQ
	if err := c.BindJSON(&req); err != nil {
		log.Println("反序列化参数失败")
		proto.Err(c, proto.BadRquest)
		return
	}

	if req.UserName == "" {
		proto.Err(c, proto.EmptyUserName)
		return
	}

	if req.PassWord == "" {
		proto.Err(c, proto.EmptyPassword)
		return
	}

	ret := services.Login(c, req.UserName, req.PassWord)
	proto.Success(c, ret)
}
