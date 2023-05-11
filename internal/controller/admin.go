package controller

import (
	"encoding/json"
	"github.com/authentication_hub/internal/domain/request"
	"github.com/authentication_hub/internal/domain/response"
	"github.com/google/uuid"
)

type AdminIer interface {
	Login([]byte) []byte
	Code([]byte) []byte
}
type AdminCtrl struct {
}

func (u *AdminCtrl) Login(message []byte) (resp []byte) {
	_ = json.Unmarshal(message, &request.LoginReq{})

	return
}
func (u *AdminCtrl) Code(message []byte) (resp []byte) {
	generateCodeReq := &request.GenerateCodeReq{}
	generateCodeResponse := &response.GenerateCodeResponse{
		BaseResp: response.BaseResp{
			Code:    "200",
			Message: "生成成功",
		},
	}
	_ = json.Unmarshal(message, generateCodeReq)
	codes := make([]string, 0)
	for i := 0; i < generateCodeReq.Count; i++ {
		uuidObj := uuid.New()
		key := uuidObj.String()
		codes = append(codes, key)
	}

	generateCodeResponse.Codes = codes
	resp, _ = json.Marshal(generateCodeResponse)
	return
}
