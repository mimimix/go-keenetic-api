package zyxel

import (
	"github.com/imroc/req/v3"
	"net/http"
)

type Zyxel struct {
	username string
	password string
	baseUrl  string
	Request  *req.Client
}

func NewZyxel(username, password, baseUrl string) *Zyxel {
	zyxel := &Zyxel{
		username: username,
		password: password,
		baseUrl:  baseUrl,
	}
	zyxel.Request = req.C().SetBaseURL(baseUrl).
		//SetCommonHeader("Content-Type", "application/json;charset=UTF-8").
		AddCommonRetryCondition(func(resp *req.Response, err error) bool {
			return resp.StatusCode == http.StatusUnauthorized
		}).
		SetCommonRetryCount(2).
		SetCommonRetryHook(zyxel.reauthManager)

	return zyxel
}

func (zyxel *Zyxel) reauthManager(resp *req.Response, _ error) {
	isSuccess, cookies := zyxel.Login()
	if isSuccess == nil {
		resp.Request.SetCookies(cookies...)
	}
}
