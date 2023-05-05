package service

import (
	"net/http"
	"time"
	"toolkit/internal/initialize"
	"toolkit/internal/logger"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/gin-gonic/gin"
)

type server interface {
	ListenAndServe() error
}

type TServerExtends struct {
	*thrift.TSimpleServer
}

func initServer(address string, router *gin.Engine) server {
	return &http.Server{
		Addr:           address,
		Handler:        router,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

func RunHttpServer() {
	Router := initialize.Routers()
	s := initServer("127.0.0.1:9100", Router)
	logger.Error(s.ListenAndServe().Error())
}

func RunThriftServer() {
	// go Httpserver()
	// go tsslServer()
	s := new(Controlloer)
	s.SetAddr("10.8.4.208", 9090)
	// if cluster.IsCluster() {
	// 	go clusterServer.ServerStart()
	// }
	s.Server()
}
