// Code generated by goctl. DO NOT EDIT.
package types

type GeneralResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type UploadResReq struct {
	Signature string `form:"signature"`
	FileName  string `form:"file_name"`
}

type UploadResResp struct {
	GeneralResponse
	Name          string `json:"name"`           // 文件名
	FileUrl       string `json:"file_url"`       // 文件的下载地址
	FileSignature string `json:"file_signature"` // 文件的标识值
	FileSize      int64  `json:"file_size"`      // 文件大小
}
