package handers

import (
	"dwn_th/proto"
	"dwn_th/services"
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

func Upload(c *gin.Context) {

	claim := Who(c)
	pwd := c.Param("pwd")
	if !strings.HasPrefix(pwd, "/") {
		pwd = "/" + pwd
	}

	file, err := c.FormFile("file")
	if err != nil {
		proto.Err(c, proto.BadRquest)
		return
	}

	f, err := file.Open()

	ret, err := services.UploadFile(c, services.UploadFileOPT{
		UserID:   int64(claim.User.ID),
		Pwd:      pwd,
		FileName: file.Filename,
		File:     f,
		Size:     file.Size,
	})
	proto.Wrap(c, ret, err)
}

func TryUpload(c *gin.Context) {

	claim := Who(c)

	var req proto.TryUploadREQ
	if err := c.BindJSON(&req); err != nil {
		log.Println("反序列化参数失败")
		proto.Err(c, proto.BadRquest)
		return
	}

	ret, err := services.TryUpload(c, req.Filehash, int64(claim.User.ID))
	proto.Wrap(c, ret, err)
}

func ConfirmUpload(c *gin.Context) {
	claim := Who(c)

	var req proto.ConfirmUploadREQ
	if err := c.BindJSON(&req); err != nil {
		log.Println("反序列化参数失败")
		proto.Err(c, proto.BadRquest)
		return
	}

	err := services.ConfirmUpload(c, req.UploadToken, int64(claim.User.ID), req.Pwd, req.FileName)
	proto.Wrap(c, nil, err)
}

// http://localhost:8080/api/file/download?user_id=1&&file_name=简历-x.pdf
func Download(c *gin.Context) {

	claim := Who(c)
	pwd := c.Param("pwd")
	if !strings.HasPrefix(pwd, "/") {
		pwd = "/" + pwd
	}

	filename, ok := c.GetQuery("file_name")
	if !ok {
		glog.Errorf("需要提供文件名")
		proto.Err(c, proto.BadRquest)
		return
	}

	file, err := services.Download(c, int64(claim.User.ID), pwd, filename)
	if err != nil {
		glog.Errorf("services.Download err: %s", err.Error())
		proto.Err(c, err)
		return
	}

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", strconv.Itoa(len(file)))
	c.Writer.Write(file)
}

func Mkdir(c *gin.Context) {

	claim := Who(c)

	var req proto.MkdirREQ
	if err := c.BindJSON(&req); err != nil {
		glog.Errorf("反序列化参数失败")
		proto.Err(c, proto.BadRquest)
		return
	}

	if strings.ContainsRune(req.Name, '/') {
		glog.Errorf("非法的目录名: %s", req.Name)
		proto.Err(c, proto.InvalidDir)
		return
	}

	err := services.Mkdir(c, req.Pwd, req.Name, int64(claim.User.ID))
	proto.Wrap(c, nil, err)
}

func List(c *gin.Context) {
	claim := Who(c)
	userId := claim.User.ID

	pwd := c.Param("pwd")
	if !strings.HasPrefix(pwd, "/") {
		pwd = "/" + pwd
	}

	ret, err := services.List(c, pwd, int64(userId))
	proto.Wrap(c, ret, err)
}

func Delete(c *gin.Context) {
	claim := Who(c)
	userId := claim.User.ID

	pwd := c.Param("pwd")
	if !strings.HasPrefix(pwd, "/") {
		pwd = "/" + pwd
	}

	filename, ok := c.GetQuery("file_name")
	if !ok {
		glog.Errorf("需要提供文件名")
		proto.Err(c, proto.BadRquest)
		return
	}

	err := services.Delete(c, pwd, int64(userId), filename)
	proto.Wrap(c, nil, err)
}
