package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// InitMiddleware 接受服务实例，并存到gin.Key中
func InitMiddleware(service []interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Keys = make(map[string]interface{})
		// 将实例存在gin.Keys中
		//c.Keys["user"] = service[0]
		c.Keys["work"] = service[0]
		c.Next()
	}
}

// ErrorMiddleware 错误处理中间件
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				c.JSON(200, gin.H{
					"code": 404,
					"msg":  fmt.Sprintf("%s", r),
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
