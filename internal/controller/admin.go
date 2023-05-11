package controller

import (
	"encoding/json"
	"github.com/authentication_hub/internal/domain/request"
	"github.com/authentication_hub/internal/domain/response"
	"github.com/authentication_hub/internal/repository/entity"
	"github.com/authentication_hub/internal/repository/persistence"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdminIer interface {
	Login([]byte) []byte
	Code([]byte) []byte
}
type AdminCtrl struct {
	db       *gorm.DB
	adminIer persistence.AdminIer
	userIer  persistence.UserIer
}

func (u *AdminCtrl) Login(message []byte) (resp []byte) {
	_ = json.Unmarshal(message, &request.AdminLoginReq{})

	return
}
func (u *AdminCtrl) Code(message []byte) (resp []byte) {
	generateCodeResponse := &response.GenerateCodeResponse{
		BaseResp: response.BaseResp{
			Code:    "200",
			Message: "生成成功",
		},
	}
	defer func() {
		resp, _ = json.Marshal(generateCodeResponse)
		return
	}()
	generateCodeReq := &request.GenerateCodeReq{}

	err := json.Unmarshal(message, generateCodeReq)
	if err != nil {
		if err != nil {
			generateCodeResponse.Code = "400"
			generateCodeResponse.Message = "请求参数错误:" + err.Error()
		}
		return
	}
	codes := make([]string, 0)
	users := make([]*entity.User, 0)
	for i := 0; i < generateCodeReq.Count; i++ {
		uuidObj := uuid.New()
		key := uuidObj.String()
		codes = append(codes, key)
		users = append(users, &entity.User{
			Mac:    "",
			IP:     "",
			Code:   key,
			Use:    0,
			Online: 0,
		})
	}
	err = u.userIer.BatchSave(u.db, users)
	if err != nil {
		generateCodeResponse.Code = "500"
		generateCodeResponse.Message = "服务器保存失败"
	}
	generateCodeResponse.Codes = codes
	return
}
