package pmas

import (
	ApiV1 "toolkit/internal/api/v1"

	"github.com/gin-gonic/gin"
)

type PMAServiceRouter struct{}

func (r *PMAServiceRouter) InitPMAServiceRouter(Router *gin.RouterGroup) {
	pmaRouter := Router.Group("PMAService")
	pmasApi := ApiV1.ApiGroupApp.PMAServiceApi
	pmaRouter.POST("sendOperTicket", pmasApi.SendOperTicket)
}
