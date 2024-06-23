package keenetic

import (
	"errors"
	"fmt"
	"github.com/imroc/req/v3"
	"net/http"
)

type Keenetic struct {
	username string
	password string
	baseUrl  string
	Request  *req.Client
}

var ErrInvalidCode = errors.New("invalid code")

func NewKeenetic(username, password, baseUrl string) *Keenetic {
	zyxel := &Keenetic{
		username: username,
		password: password,
		baseUrl:  baseUrl,
	}
	zyxel.Request = req.C().SetBaseURL(baseUrl).
		//SetCommonHeader("Content-Type", "application/json;charset=UTF-8").
		AddCommonRetryCondition(func(resp *req.Response, err error) bool {
			if !errors.Is(err, ErrInvalidCredentials) {
				return false
			}
			return resp.StatusCode == http.StatusUnauthorized
		}).
		SetCommonRetryCount(2).
		SetCommonRetryHook(func(resp *req.Response, _ error) {
			isSuccess, cookies := zyxel.Login()
			if isSuccess == nil {
				resp.Request.SetCookies(cookies...)
			}
		}).
		OnAfterResponse(func(client *req.Client, resp *req.Response) error {
			if resp.StatusCode != http.StatusOK {
				if resp.StatusCode == http.StatusUnauthorized {
					return ErrInvalidCredentials
				}
				return fmt.Errorf("%w: %d", ErrInvalidCode, resp.StatusCode)
			}
			return nil
		})

	return zyxel
}
