package main

import (
	"dwn_th/db"
	"dwn_th/handers"
	"dwn_th/model"
	"dwn_th/proto"
	"dwn_th/services"
	"dwn_th/storage"
	"encoding/base64"
	"encoding/json"
	"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func main() {

	{
		db.InitDB()
		model.Migerate()

		storage.InitOssClient()
	}

	{
		server := Route()
		server.Run()
	}

	// arr := []int{9, 4, 5, 11, 22, 8, 4, 18, 0, 2}
	// bubblesort(arr, len(arr))
	// insertsort(arr)
	// quickSort(arr, 0, len(arr)-1)
	// selectSort(arr, 10)
	// fmt.Println(arr)

}

//解析token
func Auth(c *gin.Context) {
	var claim services.Claim
	tokenStr := c.GetHeader("token")
	token, err := jwt.ParseWithClaims(tokenStr, &claim,
		func(t *jwt.Token) (interface{}, error) { return []byte("fuck"), nil })
	if err != nil {
		log.Printf("解析token失败: %s", err.Error())
		proto.Err(c, proto.InvalidToken)
		c.Abort()
		return
	}

	if claim, ok := token.Claims.(*services.Claim); ok && token.Valid {
		c.Set("user", claim)
		log.Printf("身份验证成功")
	} else {
		log.Printf("解析token失败")
		proto.Err(c, proto.InvalidToken)
		c.Abort()
		return
	}

	c.Next()
}

func Route() *gin.Engine {
	r := gin.Default()
	r.POST("/login", handers.Login)
	r.POST("/create", handers.CreateUser)

	{
		api := r.Group("/api", Auth)

		{
			user := api.Group("/user")

			user.POST("/updatepassword", handers.UpdatePsssword)
		}

		{
			file := api.Group("/file")

			file.POST("/upload/*pwd", handers.Upload)
			file.POST("/tryUpload", handers.TryUpload)
			file.POST("confirmUpload", handers.ConfirmUpload)
			file.GET("/download/*pwd", handers.Download)
			file.POST("/mkdir", handers.Mkdir)
			file.POST("/list/*pwd", handers.List)
			file.GET("/delete/*pwd", handers.Delete)

		}
	}

	return r
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

func bubblesort(s []int, m int) {
	if m == 1 {
		return
	}

	for i := 0; i < m-1; i++ {
		if s[i] > s[i+1] {
			swap(s, i, i+1)
		}
	}
	bubblesort(s, m-1)
}

func insertsort(s []int) {

	if len(s) < 2 {
		return
	}

	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			swap(s, j, j-1)
		}
	}
}

// arr := []int{9, 4, 5, 11, 22, 8, 4, 18, 0, 2}
// quickSort(arr, 0, 9)
func partition(arr []int, low, high int) int {
	pivot := arr[low]
	for low < high {
		for low < high && pivot <= arr[high] {
			high--
		}
		arr[low] = arr[high]

		for low < high && pivot >= arr[low] {
			low++
		}
		arr[high] = arr[low]
	}
	arr[low] = pivot
	return low
}
func quickSort(arr []int, low, high int) {
	if low >= high {
		return
	}
	p := partition(arr, low, high)
	quickSort(arr, low, p-1)
	quickSort(arr, p+1, high)
}

func selectSort(arr []int, n int) {
	for i := 0; i < n; i++ {
		min := i
		for j := i + 1; j < n; j++ {
			if arr[j] < arr[min] {
				min = j
			}
		}
		swap(arr, i, min)
	}
}

func swap(s []int, i, j int) {
	s[i], s[j] = s[j], s[i]
}