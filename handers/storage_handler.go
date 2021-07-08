package handers

import (
	"dwn_th/proto"
	"dwn_th/services"
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
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

	ret := services.UploadFile(c, services.UploadFileOPT{
		UserID:   int64(claim.User.ID),
		Pwd:      pwd,
		FileName: file.Filename,
		File:     f,
		Size:     file.Size,
	})
	proto.Success(c, ret)
}

func TryUpload(c *gin.Context) {

	claim := Who(c)

	var req proto.TryUpload
	if err := c.BindJSON(&req); err != nil {
		log.Println("反序列化参数失败")
		proto.Err(c, proto.BadRquest)
		return
	}

	ret := services.TryUpload(c, req.Filehash, int64(claim.User.ID))
	proto.Success(c, ret)
}

func ConfirmUpload(c *gin.Context) {
	claim := Who(c)

	var req proto.ConfirmUpload

	if err := c.BindJSON(&req); err != nil {
		log.Println("反序列化参数失败")
		proto.Err(c, proto.BadRquest)
		return
	}

	ret := services.ConfirmUpload(c, req.UploadToken, int64(claim.User.ID), req.Pwd, req.FileName)
	proto.Success(c, ret)
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
		log.Printf("需要提供文件名")
		proto.Err(c, proto.BadRquest)
		return
	}

	file, ret := services.Download(c, int64(claim.User.ID), pwd, filename)
	if ret != proto.StdSuccess {
		proto.Err(c, ret)
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
		log.Println("反序列化参数失败")
		proto.Err(c, proto.BadRquest)
		return
	}

	if strings.ContainsRune(req.Name, '/') {
		log.Printf("非法的目录名: %s", req.Name)
		proto.Err(c, proto.InvalidDir)
		return
	}

	ret := services.Mkdir(c, req.Pwd, req.Name, int64(claim.User.ID))
	proto.Success(c, ret)
}

func List(c *gin.Context) {
	claim := Who(c)
	userId := claim.User.ID

	pwd := c.Param("pwd")
	if !strings.HasPrefix(pwd, "/") {
		pwd = "/" + pwd
	}

	ret := services.List(c, pwd, int64(userId))
	proto.Success(c, ret)
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
		log.Printf("需要提供文件名")
		proto.Err(c, proto.BadRquest)
		return
	}

	ret := services.Delete(c, pwd, int64(userId), filename)
	proto.Success(c, ret)
}
