package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type wechatOption struct {
	Token         string
	AppID         string
	Secret        string
	AccessToken   string
	TokenExpoire  int64
	APITicket     string
	TicketExpoire int64

	tokenMutex  sync.Mutex
	ticketMutex sync.Mutex
}

type tokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

type ticketResp struct {
	ErrCode   int    `json:"errcode"`
	ErrMsg    string `json:"errmsg"`
	Ticket    string `json:"ticket"`
	ExpiresIn int64  `json:"expires_in"`
}

// WechatManger xxx
var WechatManger *wechatOption

func NewWechat(token string, appID string, secret string) *wechatOption {
	return &wechatOption{
		Token:         token,
		AppID:         appID,
		Secret:        secret,
		AccessToken:   "",
		TokenExpoire:  0,
		APITicket:     "",
		TicketExpoire: 0,
	}
}

func (w wechatOption) GetAccessToken() string {
	w.tokenMutex.Lock()

	if w.AccessToken == "" || w.TokenExpoire < time.Now().Unix() {
		resp, err := getToken(w.AppID, w.Secret)
		if err == nil {
			w.AccessToken = resp.AccessToken
			w.TokenExpoire = time.Now().Unix() + resp.ExpiresIn
		}
	}

	w.tokenMutex.Unlock()

	return w.AccessToken
}

func (w wechatOption) GetApiTicket() string {
	w.ticketMutex.Lock()

	if w.APITicket == "" || w.TicketExpoire < time.Now().Unix() {
		resp, err := getTicket(w.GetAccessToken())
		if err == nil {
			w.APITicket = resp.Ticket
			w.TicketExpoire = time.Now().Unix() + resp.ExpiresIn
		}
	}

	w.ticketMutex.Unlock()

	return w.APITicket
}

func getToken(appID string, secret string) (*tokenResp, error) {
	// https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=APPID&secret=APPSECRET
	response, err := http.Get(
		fmt.Sprintf(
			"https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",
			appID,
			secret,
		),
	)
	defer response.Body.Close()
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var result tokenResp
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func getTicket(accessToken string) (*ticketResp, error) {
	// https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=ACCESS_TOKEN&type=jsapi
	response, err := http.Get(
		fmt.Sprintf(
			"https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi",
			accessToken,
		),
	)
	defer response.Body.Close()
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var result ticketResp
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
