/**
 * @Author: loyeller
 * @Description:
 * @File:  ClientToken
 * @Version: 1.0.0
 * @Date: 2021/11/24 16:13
 */
package UStorageClient

import (
	"UResNode/UTool"
	"UResNode/internal/Data"
	"UResNode/internal/types"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type ClientToken struct {
	Token        string `json:"token"`
	TokenExpire  int64  `json:"token_expire"`
	RefreshToken string `json:"refresh_token"`
	RefreshUrl   string
	TokenChannel chan string
}

func NewClientToken(tk *types.ResNodeObject, refreshUrl string) *ClientToken {
	ins := &ClientToken{
		Token:        tk.Token,
		TokenExpire:  tk.TokenExpire,
		RefreshToken: tk.RefreshToken,
		RefreshUrl:   refreshUrl,
		TokenChannel: make(chan string),
	}

	go ins.FreshTokenDelay()
	return ins
}

func (ct *ClientToken) ReFreshToken() {
	ct.TokenChannel <- "refresh"
}

func (ct *ClientToken) GetTokenByRefreshToken() error {
	rfq := &types.RefreshReq{
		RefreshToken: ct.RefreshToken,
	}
	data, err := json.Marshal(rfq)
	if err != nil {
		return err
	}

	client := http.Client{Timeout: 5 * time.Second}

	req, err := http.NewRequest("POST", ct.RefreshUrl, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", ct.Token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("refresh token's resp's statuscode = %d", resp.StatusCode))
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var p types.ResNodeResp
	err = json.Unmarshal(content, &p)
	if err != nil {
		return err
	}

	if p.Code == Data.API_FAIL {
		return errors.New(fmt.Sprintf("server's resp's code = %s", p.Code))
	} else {
		ct.Token = p.Data.Token
		ct.RefreshToken = p.Data.RefreshToken
		ct.TokenExpire = p.Data.TokenExpire
	}

	return nil
}

func (ct *ClientToken) FreshTokenDelay() {
	run := true
	for ; run; {
		select {
		case cmd := <-ct.TokenChannel:
			{
				switch cmd {
				case "quit":
					{
						// 提前结束token的生命周期，不需要刷新
						run = false
					}
				case "refresh":
					{
						err := ct.GetTokenByRefreshToken()
						if err != nil {
							UTool.LogxLogFail(err.Error())
						}
					}
				}
			}
		case <-time.After(time.Duration(ct.TokenExpire/2) * time.Second):
			{
				err := ct.GetTokenByRefreshToken()
				if err != nil {
					UTool.LogxLogFail(err.Error())
				}
			}
		}
	}
	UTool.LogxFormatSuccess()
}

func (ct *ClientToken) Close() error {
	close(ct.TokenChannel)
	return nil
}
