package initialize

import (
	"toolkit/internal/router"

	"github.com/gin-gonic/gin"
)

// 初始化总路由

func Routers() *gin.Engine {
	Router := gin.Default()

	pmaRouter := router.RouterGroupApp.PMAS
	PublicGroup := Router.Group("")
	{
		// 健康监测
		PublicGroup.GET("/health", func(c *gin.Context) {
			c.JSON(200, "ok")
		})
		// TEST
		PublicGroup.GET("/test", func(c *gin.Context) {
			c.JSON(200, "ok")
		})
	}

	pmaRouter.InitPMAServiceRouter(PublicGroup)
	return Router
}
