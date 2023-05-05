package err

import (
	"runtime/debug"
	"toolkit/internal/logger"
	"toolkit/internal/service"

	"github.com/gin-gonic/gin"
)

// HandleNotFound 404
func HandleNotFound(c *gin.Context) {
	//global.NewResult(c).Error(404,"资源未找到")
	//log.Panic("资源未找到")
	service.GetResult(c).TempleError(404, "资源未找到")

}

// Recover 捕获错误
func Recover(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("panic: %v\n", r)
			debug.PrintStack() // 打印堆栈
			service.GetResult(c).Error(500000, "未知错误", errorToString(r))
			// 终止后续接口调用, 不然会继续执行
			c.Abort()
		}
	}()
	c.Next()
}

func errorToString(r interface{}) string {
	switch v := r.(type) {
	case error:
		return v.Error()
	default:
		return v.(string)
	}
}
