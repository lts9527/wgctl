package http

import (
	service "api-gateway/api/grpc/v1"
	"api-gateway/pkg/e"
	"api-gateway/pkg/res"
	"context"
	"github.com/gin-gonic/gin"
	"strings"
)

func Delete(c *gin.Context) {
	var req service.DeleteOptions
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, res.GinH(e.ERROR, err.Error(), nil))
		return
	}
	Service := c.Keys["work"].(service.ServiceClient)
	Resp, err := Service.Delete(context.Background(), &req)
	if err != nil {
		msg := strings.Replace(err.Error(), "rpc error: code = Unknown desc = ", "", -1)
		c.JSON(400, res.GinH(e.ERROR, msg, nil))
		return
	}
	r := res.Response{
		Ret:  1,
		Msg:  "请求成功",
		Data: Resp.Delete,
	}
	c.JSON(200, r)
}
