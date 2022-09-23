package res

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Ret   uint        `json:"ret,omitempty"`
	Msg   string      `json:"msg,omitempty"`
	Error string      `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

type DataList struct {
	Total uint        `json:"total,omitempty"`
	Item  interface{} `json:"item,omitempty"`
}

type TokenData struct {
	Token string      `json:"token,omitempty"`
	User  interface{} `json:"user,omitempty"`
}

func GinH(code int, msg string, data interface{}) gin.H {
	return gin.H{
		"ret":  code,
		"msg":  msg,
		"data": data,
	}
}
