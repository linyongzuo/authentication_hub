package controller

import (
	"encoding/json"
	"github.com/authentication_hub/internal/domain/response"
)

type UserIer interface {
	Heartbeat([]byte) []byte
}
type UserCtrl struct {
}

func (u *UserCtrl) Heartbeat(message []byte) (resp []byte) {
	resp, _ = json.Marshal(response.LoginResponse{
		Code:    "200",
		Message: "登陆",
	})
	return
}
