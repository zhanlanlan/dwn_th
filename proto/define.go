package proto

import "github.com/gin-gonic/gin"

type Base struct {
	Code int64
	Msg  string
	Data interface{}
}

type IErr interface {
	IntoBase() *Base
	GetCode() int64
	GetMsg() string
}

func (b *Base) GetCode() int64 {
	return b.Code
}

func (b *Base) GetMsg() string {
	return b.Msg
}

func (b *Base) IntoBase() *Base {
	return b
}

func (b *Base) Error() string {
	return b.Msg
}

func Success(c *gin.Context, value interface{}) {
	c.JSON(200, value)
}

func Err(c *gin.Context, e *Base) {
	c.JSON(200, e)
}
