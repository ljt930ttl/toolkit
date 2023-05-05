package pmas

import (
	"github.com/gin-gonic/gin"
)

type PMAServiceRouter struct{}

func (r *PMAServiceRouter) InitPMAServiceRouter(Router *gin.RouterGroup) {
	pmaRouter := Router.Group("PMAService")
	pmaRouter.POST("sendOperTicket")
}
