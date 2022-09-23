package routes

import (
	middleware "api-gateway/middlerware"
	srv "api-gateway/server/http"
	"github.com/gin-gonic/gin"
)

func NewRouter(service ...interface{}) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors(), middleware.InitMiddleware(service))
	work := r.Group("/api/v1/work/")
	{
		work.GET("ping", func(c *gin.Context) {
			c.JSON(200, "success")
		})

		work.POST("create", srv.Create)
		work.POST("show", srv.Show)
		//goland:noinspection LanguageDetectionInspection
		work.POST("ps", srv.Ps)
		work.POST("delete", srv.Delete)
		//task.POST("UpdateWgConfig", handler.UpdateWgConfig)
		//task.POST("ShowWgRunning", handler.ShowWgRunning)
		//
		//task.POST("UpdateIPAddr", handler.UpdateIPAddr)
		//task.POST("ExecPing", handler.ExecPing)
		//
		//task.POST("ShowClientLog", handler.ShowClientLog)
		//task.POST("DeleteTask", handler.UpdateScript)
	}
	return r
}
