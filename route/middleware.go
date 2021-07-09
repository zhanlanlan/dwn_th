package route

import (
	"dwn_th/proto"
	"dwn_th/services"
	"encoding/json"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

var keyFunc = func(t *jwt.Token) (interface{}, error) { return []byte("fuck"), nil }

//解析token
func Auth(c *gin.Context) {
	tokenStr := c.GetHeader("token")

	err := func() (err error) {
		if tokenStr == "" {
			glog.Warning("token不能为空")
			return proto.InvalidToken
		}

		var claim services.Claim
		token, err := jwt.ParseWithClaims(tokenStr, &claim, keyFunc)
		if err != nil {
			glog.Errorf("解析token失败: %s", err.Error())
			return proto.InvalidToken
		}

		if claim, ok := token.Claims.(*services.Claim); ok && token.Valid {
			c.Set("user", claim)
			glog.Infof("身份验证成功, user: %s", toJsonStr(claim))
		} else {
			glog.Error("解析token失败")
			return proto.InvalidToken
		}

		return
	}()
	if err != nil {
		glog.Error("auth failed err: %s", err.Error())
		c.Abort()
		return
	}

	c.Next()
}

func toJsonStr(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		glog.Errorf("toJsonStr.Marshal err: %s", err.Error())
		return ""
	}

	return string(data)
}
