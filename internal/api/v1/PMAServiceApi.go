package v1

import (
	"toolkit/internal/logger"

	"github.com/gin-gonic/gin"
)

type PMAServiceApi struct{}

func (api *PMAServiceApi) SendOperTicket(c *gin.Context) {
	logger.Debug(c.FullPath())
	logger.Debug(c.Keys)
	buf := make([]byte, 1024)
	n, _ := c.Request.Body.Read(buf)

	logger.Debug(string(buf[0:n]))
	c.JSON(200, "ok")
}

func (api *PMAServiceApi) AskAllYXAndBS(c *gin.Context) {
	logger.Debug(c.FullPath())
	logger.Debug(c.Keys)
	buf := make([]byte, 1024)
	n, _ := c.Request.Body.Read(buf)

	logger.Debug(string(buf[0:n]))
	c.JSON(200, "ok")
}
