package service

import (
	"encoding/json"
	"github.com/authentication_hub/internal/domain/request"
	"github.com/sirupsen/logrus"
)

func handlerMessage(message []byte) (resp []byte) {
	header := request.Header{}
	err := json.Unmarshal(message, &header)
	if err != nil {
		logrus.Warnf("序列化消息失败:%s", err.Error())
	}
	switch header.MessageType {
	case request.MessageHeartbeat:
		{
			return ctrl.UserIer.Heartbeat(message)
		}
	case request.MessageAdminLogin:
		{
			return ctrl.AdminCtrl().Login(message)
		}
	case request.MessageAdminGenerateCode:
		{
			return ctrl.AdminCtrl().Code(message)
		}
	}
	return
}
