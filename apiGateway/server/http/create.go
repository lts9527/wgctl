package http

import (
	service "api-gateway/api/grpc/v1"
	"api-gateway/pkg/e"
	"api-gateway/pkg/res"
	"context"
	"github.com/gin-gonic/gin"
	"strings"
)

func Create(c *gin.Context) {
	var req service.CreateOptions
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, res.GinH(e.ERROR, err.Error(), nil))
		return
	}
	//claim, _ := util.ParseToken(c.GetHeader("Authorization"))
	//req.UserID = uint32(claim.UserID)
	Service := c.Keys["work"].(service.ServiceClient)
	Resp, err := Service.Create(context.Background(), &req)
	if err != nil {
		msg := strings.Replace(err.Error(), "rpc error: code = Unknown desc = ", "", -1)
		c.JSON(400, res.GinH(e.ERROR, msg, nil))
		return
	}
	r := res.Response{
		Ret:  1,
		Msg:  "请求成功",
		Data: Resp,
	}
	c.JSON(200, r)
}
