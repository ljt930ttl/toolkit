package impl

import (
	"context"
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
	// PMA     *PMA
	Client thrift.TTransport
}

func (p *PMAImpl) RequestFunc(ctx context.Context, pmaMsg *PMAMsg) error {
	fmt.Printf("Received message: %s\n", pmaMsg)
	// response := &PMAMsg{
	// 	Head:    pmaMsg.Head,
	// 	Src:     "TestDemo",
	// 	Targets: pmaMsg.Targets,
	// 	Content: pmaMsg.Content,
	// }
	return nil
}

type NodePool struct {
	// 注册了的节点
	pool map[string]*PMANode
	// 锁
	rwLock *sync.RWMutex
}

type PMANode struct {
	IP         string
	Src        string
	Targets    []string
	Ts         thrift.TTransport
	ID         string
	IsRegister bool
	Client     *PMAServiceClient
	Sync       *sync.Mutex
}

func (np *NodePool) Register(n *PMANode) {
	logger.Debug("register.: src:%s, targets:%v, IP:%s", n.Src, n.Targets, n.IP)
	np.rwLock.Lock()
	defer np.rwLock.Unlock()
	nodeName, _ := n.GetNodeName()
	np.pool[nodeName] = n
}

func (n *PMANode) GetNodeName() (name string, er error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("GetNodeName,", err)
			logger.Error(string(debug.Stack()))
			er = errors.New("GetNodeName err")
		}
	}()
	if n.Src == "" || n.IP == "" {
		return "", errors.New("pma node err")
	}
	name = MD5(fmt.Sprint(n.Src, "PMA", n.IP))
	return
}

var NP = NodePool{
	pool:   make(map[string]*PMANode),
	rwLock: new(sync.RWMutex),
}
