package controller

import (
	"context"
	"encoding/json"
	"github.com/authentication_hub/global"
	"github.com/authentication_hub/internal/constants"
	"github.com/authentication_hub/internal/domain/request"
	"github.com/authentication_hub/internal/domain/response"
	"github.com/authentication_hub/internal/repository/entity"
	"github.com/authentication_hub/internal/repository/filter"
	"github.com/authentication_hub/internal/repository/persistence"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type UserIer interface {
	Heartbeat([]byte) []byte
	Login([]byte) []byte
	Offline([]byte) []byte
}
type UserCtrl struct {
	db      *gorm.DB
	userIer persistence.UserIer
}

func (u *UserCtrl) Heartbeat(message []byte) (resp []byte) {
	logrus.Infof("心跳:%s", string(message))
	heartbeatResp := &response.HeartbeatResp{
		BaseResp: response.BaseResp{
			Code:    "200",
			Message: "心跳成功",
		},
	}
	defer func() {
		resp, _ = json.Marshal(heartbeatResp)
		return
	}()
	req := request.HeartbeatReq{}
	err := json.Unmarshal(message, &req)
	if err != nil {
		heartbeatResp.Code = "400"
		heartbeatResp.Message = "解析请求失败"
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	global.Rdb.HSet(ctx, req.Ip, constants.KActiveTime, time.Now().Format(constants.KTimeTemplate))
	return
}
func (u *UserCtrl) Login(message []byte) (resp []byte) {
	userLoginResp := &response.UserLoginResponse{
		BaseResp: response.BaseResp{
			Code:    "200",
			Message: "用户登陆成功",
		},
	}
	defer func() {
		resp, _ = json.Marshal(userLoginResp)
		return
	}()
	req := request.UserLoginReq{}
	err := json.Unmarshal(message, &req)
	if err != nil {
		userLoginResp.Code = "400"
		userLoginResp.Message = "解析请求失败"
		return
	}
	userLoginResp.Mac = req.Mac
	user := &entity.User{}
	out, err := u.userIer.Get(u.db, user, filter.WithCode(req.Code))
	if err != nil && err == gorm.ErrRecordNotFound {
		userLoginResp.Code = "401"
		userLoginResp.Message = "认证失败，code非法"
		return
	}
	// 首次使用
	if out.Mac == "" {
		user.Online = 1
		user.Mac = req.Mac
		user.Use = 1
		user.IP = req.Ip
		_, err = u.userIer.Update(u.db, user, filter.WithCode(req.Code))
		if err != nil {
			userLoginResp.Code = "500"
			userLoginResp.Message = "服务器更新状态失败"
			return
		}
	} else {
		if out.Mac != req.Mac {
			userLoginResp.Code = "401"
			userLoginResp.Message = "服务器认证失败，mac对应错误"
			return
		}
		user.Online = 1
		user.Mac = req.Mac
		user.IP = req.Ip
		_, err = u.userIer.Update(u.db, user, filter.WithCode(req.Code))
		if err != nil {
			userLoginResp.Code = "500"
			userLoginResp.Message = "服务器更新状态失败"
			return
		}
	}
	return
}
func (u *UserCtrl) Offline(message []byte) (resp []byte) {
	req := request.UserLogoutReq{}
	err := json.Unmarshal(message, &req)
	if err != nil {
		return
	}
	user := &entity.User{IP: req.Ip}
	user.Online = 2
	_, err = u.userIer.Update(u.db, user, filter.WithIp(req.Ip))
	if err != nil {
		return
	}

	return
}
