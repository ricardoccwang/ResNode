package config

import "github.com/tal-tech/go-zero/rest"

type Config struct {
	rest.RestConf
	ResServer ResServerConf
	Node      NodeConf
	Root      string
	ResNode   ResNodeConf
}

type ResServerConf struct {
	Type       string
	Host       string
	Port       int64
	EncryptKey string
}

type NodeConf struct {
	NID           string // 节点的id
	NName         string // 节点的名称
	NType         string // server/node  这两种类型
	NHost         string // 访问节点的Host
	NPort         int64  // 访问节点的port
	NHType        string // 是http 还是 https 还是 oss 还是 obs
	NSKey         string // 访问节点需要的apitoken的key（如果有的话）
	NSSecret      string // 访问节点需要的apitoken的secret (如果有的话）
	NProviderName string // 附加信息，会显示是ali/huawei/tencent/aws/azure等
	NVersion      string // 访问节点的版本
}

type ResNodeConf struct {
	RNUrl string  // 资源的后缀
	RNHost string // 资源节点的类型
	RNPort int64 // 资源节点的类型
	RNHType string // 资源节点的类型
}