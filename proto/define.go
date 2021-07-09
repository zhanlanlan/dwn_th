package proto

import "github.com/gin-gonic/gin"

type Base struct {
	Code int64
	Msg  string
}

func (b *Base) GetCode() int64 {
	return b.Code
}

func (b *Base) GetMsg() string {
	return b.Msg
}

func (b *Base) Error() string {
	return b.Msg
}

func Success(c *gin.Context, data interface{}) {
	Wrap(c, data, nil)
}

func Err(c *gin.Context, err error) {
	Wrap(c, nil, err)
}

func Wrap(c *gin.Context, data interface{}, err error) {
	if err != nil {
		if e, ok := err.(*Base); ok {
			res(c, e.Code, e.Msg, nil)
			return
		}

		res(c, InternalErr.Code, InternalErr.Msg, nil)
		return
	}

	res(c, StdSuccess.Code, StdSuccess.Msg, data)
	return

}

func res(c *gin.Context, code int64, msg string, data interface{}) {

	res := gin.H{
		"code": code,
		"msg":  msg,
	}

	if data != nil {
		res["data"] = data
	}

	c.JSON(200, res)
}
