/**
 * @Author: loyeller
 * @Description:
 * @File:  SendFileInfoToServer
 * @Version: 1.0.0
 * @Date: 2021/11/10 15:49
 */
package UStorageClient

import (
	"UResNode/UTool"
	"UResNode/internal/Data"
	"UResNode/internal/config"
	"UResNode/internal/types"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tal-tech/go-zero/core/logx"
	"io/ioutil"
	"net/http"
	"time"
)

type ResObject struct {
	Name          string `json:"name"`           // 文件名
	FileUrl       string `json:"file_url"`       // 文件的下载地址
	FileSignature string `json:"file_signature"` // 文件的标识值
	FileSize      int64  `json:"file_size"`      // 文件大小
	Op            string `json:"op"`             // 操作方式
}

type ResReq struct {
	Ress  []*ResObject `json:"res_s"`
	Total int64       `json:"total"`
}

type ResResp struct {
	types.GeneralResponse
}
func CreateResObjectFromResUnit(downloadUrl string, unit *ResUnit) *ResObject {
	return &ResObject{
		Name:          unit.Name,
		FileUrl:       downloadUrl + "/" + unit.UrlPath,
		FileSignature: unit.FileSignature,
		FileSize:      unit.FileSize,
		Op:            unit.Op,
	}
}

func GetServerSalt(serverUrl string) (string, error) {
	hc := http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", serverUrl+"/salt", nil)
	if err != nil {
		return "", err
	}
	resp, err := hc.Do(req)
	if err != nil {
		return "", err
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	salt := types.SaltResp{}

	err = json.Unmarshal(content, &salt)
	if err != nil {
		return "", err
	}

	if salt.Code == "9999" {
		return "", errors.New("server reject node")
	}

	return salt.Data.Salt, nil
}

func SendNodeInfoToServer(serverUrl string, node config.NodeConf, encrypt string) (*types.ResNodeObject, error) {

	salt, err := GetServerSalt(serverUrl)

	if err != nil {
		return nil, err
	}

	logx.Info(fmt.Sprintf("notify node to server, nid = %s", node.NID))

	nodeQ := types.NodeReq{
		NID:           node.NID,
		NName:         node.NName,
		NType:         node.NHType,
		NHost:         node.NHost,
		NPort:         node.NPort,
		NHType:        node.NHType,
		NSKey:         node.NSKey,
		NSSecret:      node.NSSecret,
		NProviderName: node.NProviderName,
		NVersion:      node.NVersion,
		EncryptKeyJwt: UTool.SaltSecret(encrypt, salt),
	}

	hc := http.Client{Timeout: 5 * time.Second}
	data, err := json.Marshal(nodeQ)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", serverUrl, Data.HTTP_URL_NODE_LOGIN), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var p types.ResNodeResp
	err = json.Unmarshal(content, &p)
	if err != nil {
		return nil, err
	}

	if p.Code == "9999" {
		return nil, errors.New(fmt.Sprintf("notify fail, nid=%s", node.NID))
	}

	return &p.Data, nil
}

func SendFileInfoToServer(serverUrl string, downloadUrl string, unitAry []*ResUnit, token *ClientToken) error {
	logx.Info(fmt.Sprintf("send file to server, number = %d", len(unitAry)))
	resS := make([]*ResObject, 0, len(unitAry))

	for _, unit := range unitAry {
		resS = append(resS, CreateResObjectFromResUnit(downloadUrl, unit))
	}

	resQ := ResReq{
		Ress:  resS,
		Total: int64(len(resS)),
	}

	hc := http.Client{Timeout: 5 * time.Second}
	data, err := json.Marshal(resQ)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", serverUrl+"/node/res", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token.Token)
	resp, err := hc.Do(req)
	if err != nil {
		return err
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var p types.GeneralResponse
	err = json.Unmarshal(content, &p)
	if err != nil {
		return err
	}

	if p.Code == "9999" {
		return errors.New(fmt.Sprintf("notify fail, total = %d", len(unitAry)))
	}

	return nil
}
