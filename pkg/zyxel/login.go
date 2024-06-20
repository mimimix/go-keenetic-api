package zyxel

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/imroc/req/v3"
	"net/http"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type authRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (zyxel *Zyxel) Login() (error, []*http.Cookie) {
	authReq := req.C().SetBaseURL(zyxel.baseUrl)
	respTokens, errTokens := authReq.R().Get("/auth")
	if errTokens != nil {
		return errTokens, nil
	}
	xNdmChallenge := respTokens.Header.Get("X-Ndm-Challenge")
	xNdmRealm := respTokens.Header.Get("X-Ndm-Realm")
	cookies := respTokens.Cookies()
	var sessionId string
	for _, cookie := range cookies {
		if cookie.Name == "session_id" {
			sessionId = cookie.Value
		}
	}
	if sessionId == "" {
		return fmt.Errorf("session_id not found"), nil
	}

	hashermd5 := md5.New()
	hashermd5.Write([]byte(zyxel.username + ":" + xNdmRealm + ":" + zyxel.password))
	md5String := hex.EncodeToString(hashermd5.Sum(nil))
	//fmt.Println("MD5 хеш для", ":", md5String)

	input := xNdmChallenge + md5String
	hash := sha256.New()
	hash.Write([]byte(input))
	hashedString := hex.EncodeToString(hash.Sum(nil))

	zyxel.Request.ClearCookies()
	authCookies := []*http.Cookie{
		{Name: "session_id", Value: sessionId},
		{Name: "_authorized", Value: zyxel.username},
		{Name: "sysmode", Value: "router"},
	}
	zyxel.Request.SetCommonCookies(authCookies...)

	fmt.Println(sessionId)

	authResp, errResp := authReq.R().
		SetBody(&authRequest{Login: zyxel.username, Password: hashedString}).
		Post("/auth")
	if errResp != nil {
		return errResp, nil
	}
	switch {
	case authResp.StatusCode == http.StatusOK:
		return nil, authCookies
	case authResp.StatusCode == http.StatusUnauthorized:
		return ErrInvalidCredentials, nil
	default:
		return fmt.Errorf("unexpected status code: %d", authResp.StatusCode), nil
	}
}
