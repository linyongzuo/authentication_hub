package response

import "github.com/authentication_hub/internal/domain/request"

type BaseResp struct {
	Code        string              `json:"code"`        // 登陆返回
	Message     string              `json:"message"`     // 登陆返回错误码
	MessageType request.MessageType `json:"messageType"` // 消息类型
}
