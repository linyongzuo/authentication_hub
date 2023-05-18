package controller

import (
	"context"
	"encoding/json"
	"github.com/authentication_hub/global"
	"github.com/authentication_hub/internal/domain/request"
	"github.com/authentication_hub/internal/domain/response"
	"github.com/authentication_hub/internal/repository/entity"
	"github.com/authentication_hub/internal/repository/filter"
	"github.com/authentication_hub/internal/repository/persistence"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type AdminIer interface {
	Login([]byte) []byte
	Code([]byte) []byte
	UserInfo([]byte) []byte
}
type AdminCtrl struct {
	db       *gorm.DB
	adminIer persistence.AdminIer
	userIer  persistence.UserIer
}

func (u *AdminCtrl) Login(message []byte) (resp []byte) {
	adminLoginResponse := &response.AdminLoginResponse{
		BaseResp: response.BaseResp{
			Code:        "200",
			Message:     "登陆成功",
			MessageType: request.MessageAdminLogin,
		},
		UserName: "",
	}
	adminLoginReq := &request.AdminLoginReq{}
	defer func() {
		resp, _ = json.Marshal(adminLoginResponse)
		return
	}()
	err := json.Unmarshal(message, adminLoginReq)
	if err != nil {
		if err != nil {
			adminLoginResponse.Code = "400"
			adminLoginResponse.Message = "请求参数错误:" + err.Error()
		}
		return
	}
	adminLoginResponse.UserName = adminLoginReq.UserName
	_, err = u.adminIer.Get(u.db, &entity.Admin{
		UserName: adminLoginReq.UserName,
		Password: adminLoginReq.Password,
	})
	if err != nil {
		adminLoginResponse.Code = "404"
		adminLoginResponse.Message = "用户民密码错误:" + err.Error()
	}
	adminLoginResponse.UserName = adminLoginReq.UserName
	return
}
func (u *AdminCtrl) Code(message []byte) (resp []byte) {

	generateCodeResponse := &response.GenerateCodeResponse{
		BaseResp: response.BaseResp{
			Code:        "200",
			Message:     "生成成功",
			MessageType: request.MessageAdminGenerateCode,
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
func (u *AdminCtrl) UserInfo(message []byte) (resp []byte) {
	userInfo := &response.UserInfoResponse{
		BaseResp: response.BaseResp{
			Code:        "200",
			Message:     "获取成功",
			MessageType: request.MessageUserInfo,
		},
	}
	defer func() {
		resp, _ = json.Marshal(userInfo)
		return
	}()
	count, err := u.userIer.Count(u.db, &entity.User{}, filter.WithUsed())
	if err != nil {
		userInfo.Code = "400"
		userInfo.Message = "获取注册人数失败:" + err.Error()
	}
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	keys := global.Rdb.Keys(ctx, "*").Val()
	userInfo.OnlineCount = len(keys)
	userInfo.RegisterCount = int(count)
	userInfo.OfflineCount = userInfo.RegisterCount - userInfo.OnlineCount
	return
}
