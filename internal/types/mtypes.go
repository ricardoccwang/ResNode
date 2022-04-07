/**
 * @Author: loyeller
 * @Description:
 * @File:  mtypes
 * @Version: 1.0.0
 * @Date: 2021/11/12 9:00
 */
package types

type NodeReq struct {
	NID           string `json:"n_id"`            // 节点的id
	NName         string `json:"n_name"`          // 节点的名称
	NType         string `json:"n_type"`          // server/node  这两种类型
	NHost         string `json:"n_host"`          // 访问节点的host
	NPort         int64  `json:"n_port"`          // 访问节点的port
	NHType        string `json:"n_htype"`         // 是http 还是 https 还是 oss 还是 obs
	NSKey         string `json:"n_secret_key"`    // 访问节点需要的apitoken的key（如果有的话）
	NSSecret      string `json:"n_secret_secret"` // 访问节点需要的apitoken的secret (如果有的话）
	NProviderName string `json:"n_provider_name"` // 附加信息，会显示是ali/huawei/tencent/aws/azure等
	NVersion      string `json:"n_version"`       // 访问节点的版本
	EncryptKeyJwt string `json:"encrypt_key_jwt"` // Jwt格式的访问秘钥
}

type SaltObject struct {
	Salt string `json:"salt"`
}

type SaltResp struct {
	GeneralResponse
	Data SaltObject `json:"data,omitempty"`
}


type ResNodeObject struct {
	Token        string `json:"token"`
	TokenExpire  int64  `json:"token_expire"`
	RefreshToken string `json:"refresh_token"`
}

type ResNodeResp struct {
	GeneralResponse
	Data ResNodeObject `json:"data,omitempty"`
}

type RefreshReq struct {
	RefreshToken string `json:"refresh_token"`
}