package route

import (
	"dwn_th/handers"

	"github.com/gin-gonic/gin"
)

func Route(r *gin.Engine) {
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

}
