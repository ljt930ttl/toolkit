package service

import (
	"net/http"
	"toolkit/internal/logger"

	"github.com/gin-gonic/gin"
)

type Result struct {
	*gin.Context
}

// ResultCont 返回的结果：
type ResultCont struct {
	Code int         `json:"code"` //提示代码
	Msg  string      `json:"msg"`  //提示信息
	Err  string      `json:"err"`  //错误信息
	Data interface{} `json:"data"` //数据
}

func GetResult(c *gin.Context) *Result {
	return &Result{Context: c}
}

// Success 成功
func (r *Result) Success(data interface{}) {
	//if data == nil {
	//	data = gin.H{}
	//}
	res := ResultCont{}
	res.Code = 0
	res.Data = data
	logger.Info("成功:", res)
	r.JSON(http.StatusOK, res)
}
func (r *Result) TempleSuccess(data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	r.HTML(http.StatusOK, "index.html", data)
}

// Error 错误
func (r *Result) Error(code int, msg, err string) {
	res := ResultCont{}
	res.Code = code
	res.Msg = msg
	res.Err = err
	//res.Data = gin.H{}
	logger.Error("错误:", res)
	r.JSON(http.StatusOK, res)
}

func (r *Result) TempleError(code int, msg string) {
	res := ResultCont{}
	res.Code = code
	res.Msg = ""
	res.Err = ""
	res.Data = gin.H{
		"err":      msg,
		"rootPath": "/",
	}
	r.HTML(http.StatusOK, "err.html", res)
}
