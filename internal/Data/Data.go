/**
 * @Author: loyeller
 * @Description:
 * @File:  Data
 * @Version: 1.0.0
 * @Date: 2021/11/15 12:00
 */
package Data

const (
	MaxSendFileInfoCount = 1000
)

const (
	SickCount = 5 // 5分钟
)

const (
	SyncNodeOperatorAdd = "add"
	SyncNodeOperatorDelete = "delete"
	SyncNodeOperatorUpdate = "update"
)

const (
	API_SUCCESS = "0000"
	API_FAIL = "9999"
)

const (
	LOGX_SUCCESS_FMT = "[%s] IS SUCCESS"
	LOGX_FAIL_FMT = "[%s] IS FAIL;[Result] %s"
	LOGX_WARNING_FMT = "[%s] IS WARNING;[Reason] %s"
)

const (
	HTTP_URL_NODE_LOGIN = "/node"
	HTTP_URL_REFRESH_TOKEN = "/node/refreshtoken"
)