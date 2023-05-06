package handle

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"runtime/debug"
	"strings"
	"sync"
	"toolkit/internal/logger"
	. "toolkit/internal/protocol/gen-go/PMA"
	. "toolkit/internal/utils"

	"github.com/apache/thrift/lib/go/thrift"
)

type NodePool struct {
	pool map[*PMAServiceNode]bool
	// 注册了的节点
	poolNode map[string]*PMAServiceNode
	// 锁
	rwLock *sync.RWMutex
}

type PMAServiceNode struct {
	Addr       string
	Src        string
	IP         string
	Targets    []string
	msg        string
	ID         string
	IsRegister bool
	Ts         thrift.TTransport
	Client     *PMAServiceClient
	Sync       *sync.Mutex
}

func (np *NodePool) AddConn(n *PMAServiceNode) {
	logger.Debug("AddConn...: ")
	np.rwLock.Lock()
	defer np.rwLock.Unlock()
	np.pool[n] = true
}

func (np *NodePool) LenPool() int {
	return len(np.pool)
}

func (n *PMAServiceNode) Register(clinet thrift.TTransport, msg *PMAMsg) {
	n.Sync.Lock()
	defer n.Sync.Unlock()
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Register,", err)
			logger.Error(string(debug.Stack()))
			n.msg = "Register err, 内部错误!"
		}
	}()
	logger.Info("Received message-register:\n", msg.Content)
	addr := new(Addr)
	json.Unmarshal([]byte(msg.Content), &addr)

	n.Addr = clinet.(*thrift.TSocket).Addr().String()
	n.Src = msg.Src
	n.Targets = msg.Targets
	n.IP = addr.IP

	NP.poolNode[n.Src] = n
	n.IsRegister = true
}

func (n *PMAServiceNode) RegisterAck(clinet thrift.TTransport, msg *PMAMsg) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Register,", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	ackMsg := &PMAMsg{
		Head: make(map[string]string),
	}
	ackMsg.Head["func"] = "register"
	ackMsg.Head["time"] = NowTime()
	ackMsg.Head["requestId"] = msg.Head["requestId"]
	if n.IsRegister {
		ackMsg.Head["returnCode"] = "0"
		ackMsg.Head["returnMsg"] = ""
	} else {
		ackMsg.Head["returnCode"] = "-2"
		ackMsg.Head["returnMsg"] = n.msg
	}

	ackMsg.Head["frame"] = "1/1"

	ackMsg.Src = "cyg.test"
	ackMsg.Targets = append(ackMsg.Targets, msg.Src)
	ackMsg.Content = msg.Content
	logger.Debug("register msg ack:\n", ackMsg)
	n.Client.RequestFunc(context.Background(), ackMsg)

}

func (n *PMAServiceNode) ReceSendMsg(clinet thrift.TTransport, msg *PMAMsg) {
	n.Sync.Lock()
	defer n.Sync.Unlock()
	defer func() {
		if err := recover(); err != nil {
			logger.Error("SendMsg,", err)
			logger.Error(string(debug.Stack()))
			n.msg = "SendMsg err, 内部错误!"
		}
	}()

	logger.Info("Received message-send:\n", msg.Content)
	if frame, ok := msg.Head["frame"]; ok {
		arr := strings.Split(frame, "/")
		if len(arr) == 2 {
			ContentTemp += msg.Content
			if arr[0] == arr[1] {
				ContentTotal = ContentTemp
				ContentTemp = ""
			} else {
				return
			}
		} else {
			n.msg = "SendMsg, frame syntax"
			return
		}
	} else {
		ContentTotal = msg.Content
	}

}

func (n *PMAServiceNode) SendMsg(requestId, target, ContentTotal string) {
	ackMsg := &PMAMsg{
		Head: make(map[string]string),
	}
	ackMsg.Head["func"] = "send"
	ackMsg.Head["time"] = NowTime()
	ackMsg.Head["requestId"] = requestId
	ackMsg.Head["returnCode"] = "1"
	ackMsg.Head["returnMsg"] = ""

	ackMsg.Head["frame"] = "1/1"
	ackMsg.Src = "cyg.test"
	ackMsg.Targets = append(ackMsg.Targets, target)

	// 生成content
	m := make(map[string]interface{}, 1)
	err := json.Unmarshal([]byte(ContentTotal), &m) //第二个参数要地址传递
	if err != nil {
		logger.Error("err = ", err)
		return
	}
	method := m["method"].(string)
	if handle, ok := processMap[method]; ok {
		// logger.Debug("found function", &function)
		ackMsg.Content = handle.Process(ContentTotal)
	} else {
		ackMsg.Head["returnCode"] = "-2"
		ackMsg.Head["returnMsg"] = fmt.Sprintf("not found function method:%s", method)
	}

	logger.Debug("Send msg ack:\n", ackMsg)
	n.Client.RequestFunc(context.Background(), ackMsg)

}

func GetNodeName(src, ip string) (name string, er error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("GetNodeName,", err)
			logger.Error(string(debug.Stack()))
			er = errors.New("GetNodeName err")
		}
	}()
	if src == "" || ip == "" {
		return "", errors.New("pma node err, src or ip is null")
	}
	name = MD5(fmt.Sprint(src, "PMA", ip))
	return
}

var NP = NodePool{
	pool:     make(map[*PMAServiceNode]bool),
	poolNode: make(map[string]*PMAServiceNode),
	rwLock:   new(sync.RWMutex),
}

var ContentTotal = ""
var ContentTemp = ""
