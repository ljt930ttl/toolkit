package handle

import (
	"context"
	"errors"
	"toolkit/internal/logger"

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
	// logger.Debug("Received message: ", pmaMsg)
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
		logger.Debug("not found node......", NP.poolNode, pmaMsg.Src)
	}
}
