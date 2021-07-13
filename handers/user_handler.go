package handers

import (
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
		glog.Warning("非法的用户名 UserName: %s", req.UserName)
		proto.Err(c, proto.BadUserName)
		return
	}

	if req.PassWord == "" {
		glog.Warning("空的密码 UserName: %s", req.UserName)
		proto.Err(c, proto.BadPassword)
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
		glog.Error("NewPassWord can not be empty")
		proto.Err(c, proto.BadRquest)
		return
	}

	err := services.UpdatePsssword(c, claim.User.UserName, req.NewPassWord)
	proto.Wrap(c, nil, err)
}

func Login(c *gin.Context) {
	var req proto.LoginREQ
	if err := c.BindJSON(&req); err != nil {
		glog.Errorf("反序列化参数失败 err: %s", err.Error())
		proto.Err(c, proto.BadRquest)
		return
	}

	if req.UserName == "" {
		glog.Error("Login, UserName can not be empty")
		proto.Err(c, proto.EmptyUserName)
		return
	}

	if req.PassWord == "" {
		glog.Error("Login, PassWord can not be empty")
		proto.Err(c, proto.EmptyPassword)
		return
	}

	ret, err := services.Login(c, req.UserName, req.PassWord)
	proto.Wrap(c, ret, err)
}
