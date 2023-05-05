package v1

import (
	"toolkit/internal/logger"

	"github.com/gin-gonic/gin"
)

func SendOperTicket(c *gin.Context) {
	logger.Debug(c.Params)
	c.JSON(200, "ok")
}
