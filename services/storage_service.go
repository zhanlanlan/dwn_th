package services

import (
	"context"

	"crypto/sha1"
	"crypto/sha256"
	"dwn_th/model"
	"dwn_th/proto"
	"dwn_th/storage"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"

	"path"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/golang/glog"
)

func utilGenerateKey(fileName string) string {
	sum := sha256.Sum256([]byte(fileName + strconv.FormatInt(time.Now().UnixNano(), 10)))

	return base64.StdEncoding.EncodeToString(sum[:])
}

func createFileKey() string {
	now := time.Now().UnixNano()
	rd := rand.Int63n(1000)

	return strconv.FormatInt(now+rd, 16)
}

func Download(ctx context.Context, userID int64, pwd string, fileName string) (data []byte, err error) {

	dbF, err := model.GetFileByOwner(ctx, userID, pwd, fileName)
	if errors.Is(err, proto.FileNotFound) {
		glog.Errorf("model.GetFileByOwner err: %s", err.Error())
		err = proto.FileNotFound
		return
	} else if err != nil {
		glog.Errorf("model.GetFileByOwner err: %s", err.Error())
		return
	}

	dbfile, err := model.GetFileById(ctx, dbF.FileId)
	if err != nil {
		glog.Errorf("model.GetFileById err: %s", err.Error())
		return
	}

	rdcloser, err := storage.Get(dbfile.URI)
	if err != nil {
		glog.Errorf("storage.Get err: %s", err.Error())
		return
	}
	defer rdcloser.Close()

	data, err = ioutil.ReadAll(rdcloser)
	if err != nil {
		glog.Errorf("下载文件出错: %s", err.Error())
		return
	}

	return
}

func Mkdir(ctx context.Context, pwd, name string, userid int64) (err error) {
	// 检查父目录是否存在
	pwd = path.Clean(pwd)

	// 对根目录做特殊处理  根目录默认存在
	if pwd != "/" {
		fatherdir, fathername := path.Split(pwd)
		fatherdir = path.Clean(fatherdir)
		_, err = model.GetDir(ctx, fatherdir, fathername)
		if errors.Is(err, proto.DirNotExist) {
			err = proto.DirNotExist
			return
		}

	}

	// 校验文件是否已存在
	_, err = model.GetUserFileByPwdANDFileName(ctx, pwd, name)
	if errors.Is(err, proto.FileNotFound) {
		// 文件不存在  可以创建目录
		uf := model.UserFile{

			UserId:   userid,
			Pwd:      pwd,
			FileName: name,
			FileType: model.TypeDir,
		}
		err := model.CreateUserFile(ctx, &uf)
		if err != nil {
			glog.Errorf("model.CreateUserFile err: %s", err.Error())
			return proto.InternalErr
		}

	} else if err != nil {
		// 出现其他错误
		glog.Errorf("mkdir err: %s", err.Error())
		return proto.InternalErr
	} else {
		// 文件已存在 不能创建
		glog.Errorf("mkdir failed: file already exist")
		return proto.FileAlreadyExist
	}

	return
}

func List(ctx context.Context, pwd string, userId int64) (data []proto.ListRES, err error) {

	userfiles, err := model.GetFileList(ctx, pwd, userId)
	if err != nil {
		glog.Errorf("model.GetFileList err: %s", err.Error())
		return
	}

	listRes := make([]proto.ListRES, len(userfiles))
	for i, userfile := range userfiles {
		list := proto.ListRES{
			Pwd:  userfile.Pwd,
			Name: userfile.FileName,
			Ext:  userfile.Ext,
			Type: int64(userfile.FileType),
		}
		listRes[i] = list
	}

	data = listRes
	return
}

func Delete(ctx context.Context, pwd string, userID int64, fileName string) (err error) {
	err = model.DeleteFile(ctx, userID, pwd, fileName)
	if err != nil {
		glog.Errorf("model.DeleteFile err: %s", err.Error())
		return
	}

	return
}

func TryUpload(ctx context.Context, filehash string, userid int64) (data proto.TryUploadRES, err error) {
	file, err := model.GetFileByHash(ctx, filehash)
	if err != nil {
		glog.Errorf("model.GetFileByHash err: %s", err.Error())
		err = proto.FileNotFound
		return
	}

	upload := UploadToken{
		FileKey:    file.FileKey,
		ExpireTime: 3600,
		UserID:     userid,
	}
	uploadtoken, err := GetToken(upload)
	if err != nil {
		glog.Errorf("GetToken err: %s", err.Error())
		return
	}

	data.Token = uploadtoken
	return
}

type UploadToken struct {
	FileKey    string
	ExpireTime int64
	UserID     int64
}

func (u *UploadToken) Marshal() string {
	data, _ := json.Marshal(u)
	return base64.StdEncoding.EncodeToString(data)
}

func (u *UploadToken) Unmarshal(token string) error {
	data, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, u)
}

type UploadTokenClaim struct {
	Uploadtoken UploadToken
	jwt.StandardClaims
}

func GetToken(upt UploadToken) (token string, err error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, UploadTokenClaim{
		Uploadtoken: upt,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 6).Unix(),
		},
	})

	return t.SignedString([]byte("fuck"))
}
func VerifyToken(t string) *UploadTokenClaim {

	var uptclaim UploadTokenClaim
	token, err := jwt.ParseWithClaims(t, &uptclaim, func(token *jwt.Token) (interface{}, error) { return []byte("fuck"), nil })

	if err != nil {
		log.Printf("解析token失败: %s", err.Error())
		return nil
	}

	if uptclaim, ok := token.Claims.(*UploadTokenClaim); ok && token.Valid {
		log.Printf("身份验证成功")
		return uptclaim
	} else {
		log.Printf("解析token失败")

		return nil
	}
}

type UploadFileOPT struct {
	FileName string
	UserID   int64
	Pwd      string
	File     io.Reader
	Size     int64
}

func UploadFile(ctx context.Context, opt UploadFileOPT) (data proto.UploadFileRES, err error) {

	pwd := path.Clean(opt.Pwd)
	fmt.Println("pwd: ", pwd)
	if pwd != "/" {
		fadir, faname := path.Split(pwd)
		fadir = path.Clean(fadir)
		_, err := model.GetDir(ctx, fadir, faname)
		if errors.Is(err, proto.DirNotExist) {
			err = proto.DirNotExist
		}
	}
	upt := UploadToken{
		FileKey:    createFileKey(),
		ExpireTime: 3600,
		UserID:     opt.UserID,
	}

	_, err = model.GetFile(ctx, pwd, opt.FileName)
	if errors.Is(err, proto.FileNotFound) {
		objName := utilGenerateKey(opt.FileName)
		err = storage.Put(objName, opt.File)
		if err != nil {
			glog.Errorf("storage.Put err: %s", err.Error())
			return
		}

		// 更新数据库 关联用户和文件
		f := model.File{
			FileKey: upt.FileKey,
			URI:     objName,
			Size:    opt.Size,
			Hash:    hashFile(opt.File),
		}

		err = model.CreateFile(ctx, &f)
		if err != nil {
			glog.Errorf("model.CreateFile err: %s", err.Error())
			return
		}
	} else if err != nil {
		// 出现其他错误
		glog.Errorf("model.GetFile err: %s", err.Error())
		return
	} else {
		// 文件已存在 不能创建
		glog.Errorf("UploadFile failed: file already exist")
		err = proto.FileAlreadyExist
		return
	}

	// 上传到oss

	uploadtoken, err := GetToken(upt)
	if err != nil {
		glog.Errorf("GetToken err: %s", err.Error())
		err = proto.GetTokenFail
		return
	}

	data.Token = uploadtoken
	return
}

func hashFile(f io.Reader) string {
	filebyte, err := ioutil.ReadAll(f)
	if err != nil {
		log.Printf("read file failed")
	}

	h := sha1.New()

	_, err = h.Write(filebyte)

	if err != nil {
		log.Printf("write filebyte array failed")
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}

func ConfirmUpload(ctx context.Context, uploadtoken string, userid int64, pwd string, filename string) (err error) {

	upt := VerifyToken(uploadtoken)

	file, err := model.GetFileByFileKey(ctx, upt.Uploadtoken.FileKey)
	if err != nil {
		glog.Errorf("model.GetFileByFileKey err: %s", err.Error())
		return proto.FileNotFound
	}

	uf := model.UserFile{

		UserId:   userid,
		FileId:   int64(file.ID),
		Pwd:      pwd,
		FileName: filename,
		FileType: 1,
		Ext:      "",
	}
	err = model.CreateUserFile(ctx, &uf)
	if err != nil {
		glog.Errorf("model.CreateUserFile err: %s", err.Error())
		return
	}
	return

}

type ShareToken struct {
	UserID     int64
	ExpireTime int64
	FileName   string
	FileExt    string
	FileKey    string
}

type ShareTokenClaim struct {
	Sharetoken ShareToken
	jwt.StandardClaims
}

func GetShareToken(st ShareToken) (token string, err error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, ShareTokenClaim{
		Sharetoken: st,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 6).Unix(),
		},
	})

	return t.SignedString([]byte("fuck"))
}
func VerifyShareToken(t string) *ShareTokenClaim {

	var stclaim ShareTokenClaim
	token, err := jwt.ParseWithClaims(t, &stclaim, func(token *jwt.Token) (interface{}, error) { return []byte("fuck"), nil })

	if err != nil {
		log.Printf("解析token失败: %s", err.Error())
		return nil
	}

	if stclaim, ok := token.Claims.(*ShareTokenClaim); ok && token.Valid {
		log.Printf("身份验证成功")
		return stclaim
	} else {
		log.Printf("解析token失败")

		return nil
	}
}
