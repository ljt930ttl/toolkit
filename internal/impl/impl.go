package impl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"runtime/debug"
	"sync"
	"toolkit/internal/logger"
	. "toolkit/internal/protocol/gen-go/PMA"
	. "toolkit/internal/utils"

	"github.com/apache/thrift/lib/go/thrift"
)

type PMAImpl struct {
	Ip      string
	Port    int
	Account string
	// Pub    string //发布id
	Node   *PMAServiceNode
	Client thrift.TTransport
}
type Addr struct {
	IP string `json:"ip"`
}

func (p *PMAImpl) RequestFunc(ctx context.Context, pmaMsg *PMAMsg) error {
	fmt.Printf("Received message: %s\n", pmaMsg)
	if p.Node == nil {
		return errors.New("内部错误,连接失败,node为空")
	}
	if pmaMsg.Head["func"] == "register" {
		go p.register(ctx, pmaMsg)
	}

	if pmaMsg.Head["func"] == "send" {

		go p.sendMsg(ctx, pmaMsg)
	}

	return nil
}

func (p *PMAImpl) register(ctx context.Context, pmaMsg *PMAMsg) {
	p.Node.Register(p.Client, pmaMsg)
	p.Node.RegisterAck(p.Client, pmaMsg)
}

func (p *PMAImpl) sendMsg(ctx context.Context, pmaMsg *PMAMsg) {
	if node, ok := NP.poolNode[pmaMsg.Src]; ok {
		node.SendMsg(p.Client, pmaMsg)
	} else {

	}
}

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
	logger.Debug("AddConn.: src:%s, targets:%v, IP:%s", n.Src, n.Targets, n.Addr)
	np.rwLock.Lock()
	defer np.rwLock.Unlock()
	np.pool[n] = true
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
	addr := new(Addr)
	json.Unmarshal([]byte(msg.Content), &addr)

	n.Addr = clinet.(*thrift.TSocket).Addr().String()
	n.Src = msg.Src
	n.Targets = msg.Targets
	n.IP = addr.IP
	// thrift.TServerSocket(p.Client)
	// name, err := GetNodeName(n.Src, n.IP)
	// if err != nil {
	// 	n.msg = err.Error()
	// }

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
	ackMsg := new(PMAMsg)
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
	n.Client.RequestFunc(context.Background(), ackMsg)

}

func (n *PMAServiceNode) SendMsg(clinet thrift.TTransport, msg *PMAMsg) {
	n.Sync.Lock()
	defer n.Sync.Unlock()
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Register,", err)
			logger.Error(string(debug.Stack()))
			n.msg = "Register err, 内部错误!"
		}
	}()
	m := make(map[string]interface{}, 1)
	ackMsg := new(PMAMsg)
	ackMsg.Head["func"] = "send"
	ackMsg.Head["time"] = NowTime()
	ackMsg.Head["requestId"] = msg.Head["requestId"]

	ackMsg.Head["frame"] = "1/1"
	ackMsg.Src = "cyg.test"
	ackMsg.Targets = append(ackMsg.Targets, msg.Src)

	err := json.Unmarshal([]byte(msg.Content), &m) //第二个参数要地址传递
	if err != nil {
		logger.Error("err = ", err)
		return
	}
	method := m["method"].(string)
	ackMsg.Content = processMap[method].Process(msg.Content)
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
	poolNode: make(map[string]*PMAServiceNode),
	rwLock:   new(sync.RWMutex),
}
