package service

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	. "toolkit/internal/impl"
	"toolkit/internal/logger"
	. "toolkit/internal/protocol/gen-go/PMA"

	"github.com/apache/thrift/lib/go/thrift"
)

type Controlloer struct {
	Port int    `json:"port"`
	Ip   string `json:"ip"`
}

var cfg *thrift.TConfiguration = nil

func ServerStart() {
	// go Httpserver()
	// go tsslServer()
	s := new(Controlloer)
	s.SetAddr("127.0.0.1", 9090)
	// if cluster.IsCluster() {
	// 	go clusterServer.ServerStart()
	// }
	s.Server()
}

func (t *Controlloer) SetAddr(ip string, port int) {
	t.Port = port
	t.Ip = ip
}

func (t *Controlloer) ListenAddr() string {
	return fmt.Sprint(t.Ip, ":", t.Port)
}

func (t *Controlloer) Server() {
	transportFactory := thrift.NewTBufferedTransportFactory(1024)
	protocolFactory := thrift.NewTBinaryProtocolFactoryConf(cfg)
	serverTransport, err := thrift.NewTServerSocket(t.ListenAddr())
	if err != nil {
		logger.Error("server:", err.Error())
		panic(err.Error())
	}
	handler := new(PMAImpl)
	processor := NewPMAServiceProcessor(handler)
	server := NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
	fmt.Println("server listen:", t.ListenAddr())
	Listen(server, 100)
	if err == nil {
		for {
			client, err := Accept(server)
			if err == nil {
				go controllerHandler(client)
			}
		}
	}
}

func Listen(server *TSimpleServer, count int) (err error) {
	if count <= 0 {
		err = errors.New("")
		return
	}
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Listen,", err)
			logger.Error(string(debug.Stack()))
			count--
			Listen(server, count)
		}
	}()
	err = server.Listen()
	return
}

func Accept(server *TSimpleServer) (client thrift.TTransport, err error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Accept,", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	client, err = server.Accept()
	return
}

func controllerHandler(tt thrift.TTransport) {
	isclose := false
	var gorutineclose *bool = &isclose
	defer func() {
		if err := recover(); err != nil {
			logger.Error("controllerHandler,", err)
			*gorutineclose = true
		}
	}()

	defer func() { tt.Close() }()
	monitorChan := make(chan string, 1)

	go PMAProcessor(tt, gorutineclose, monitorChan)
	<-monitorChan
	//	errormsg := <-monitorChan
	//	logger.Error("errormsg:", errormsg)
}

func NewPMAClient(tt thrift.TTransport) *PMAServiceClient {
	transportFactory := thrift.NewTBufferedTransportFactory(1024)
	protocolFactory := thrift.NewTBinaryProtocolFactoryConf(cfg)
	useTransport, _ := transportFactory.GetTransport(tt)
	return NewPMAServiceClientFactory(useTransport, protocolFactory)
}

func PMAProcessor(client thrift.TTransport, gorutineclose *bool, monitorChan chan string) error {
	defer func() {
		if err := recover(); err != nil {
			//			logger.Error(string(debug.Stack()))
			logger.Warn("processor:", err)
		}
	}()
	defer func() {
		if err := recover(); err != nil {
			logger.Warn("processor:", err)
		}
		*gorutineclose = true
		monitorChan <- "Processor end"
	}()
	node := &PMAServiceNode{
		Ts:     client,
		Client: NewPMAClient(client),
	}
	NP.AddConn(node)

	protocol := thrift.NewTBinaryProtocolConf(client, cfg)
	handler := &PMAImpl{Client: client, Node: node}
	processor := NewPMAServiceProcessor(handler)
	for {
		ok, err := processor.Process(context.Background(), protocol, protocol)
		if err, ok := err.(thrift.TTransportException); ok && err.TypeId() == thrift.END_OF_FILE {
			return nil
		} else if err != nil {
			return err
		}
		if !ok {
			logger.Error("Processor error:", err)
			break
		}
	}
	return nil
}
