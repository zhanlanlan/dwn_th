package handers

import (
	"log"

	"dwn_th/proto"
	"dwn_th/services"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var req proto.CreateUserREQ
	if err := c.BindJSON(&req); err != nil {
		log.Println("反序列化参数失败")
		proto.Err(c, proto.BadRquest)
		return
	}

	if req.UserName == "" || len([]rune(req.UserName)) > 25 {
		proto.Err(c, proto.BadUserName)
		return
	}

	if req.PassWord == "" {
		proto.Err(c, proto.BadPassword)
		return
	}

	ret := services.CreateUser(c, req)
	proto.Success(c, ret)
}

func UpdatePsssword(c *gin.Context) {
	var req proto.UpdatePssswordREQ
	if err := c.BindJSON(&req); err != nil {
		log.Println("反序列化参数失败")
		proto.Err(c, proto.BadRquest)
		return
	}

	claim := Who(c)

	if req.NewPassWord == "" {
		//
		return
	}

	ret := services.UpdatePsssword(c, claim.User.UserName, req.NewPassWord)
	proto.Success(c, ret)
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
