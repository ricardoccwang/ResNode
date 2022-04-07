package svc

import (
	"UResNode/UStorageClient"
	"UResNode/internal/config"
	"fmt"
	"github.com/tal-tech/go-zero/core/logx"
)

type ServiceContext struct {
	Config config.Config
	Client *UStorageClient.ClientManager
}

func NewServiceContext(c config.Config) *ServiceContext {
	m := UStorageClient.NewClientManager(c.Root,
		fmt.Sprintf("%s://%s:%d", c.ResServer.Type, c.ResServer.Host, c.ResServer.Port),
		fmt.Sprintf("%s://%s:%d/%s", c.ResNode.RNHType, c.ResNode.RNHost, c.ResNode.RNPort, c.ResNode.RNUrl),
		c.Node,
		c.ResServer.EncryptKey)
	err := m.StartClient()
	if err != nil {
		logx.Severe(fmt.Sprintf("[StartClient] IS FAIL;[Result] %s", err.Error()))
		panic(fmt.Sprintf("err = %s", err.Error()))
	}
	return &ServiceContext{
		Config: c,
		Client: m,
	}
}
