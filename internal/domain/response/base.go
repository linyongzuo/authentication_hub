package response

type BaseResp struct {
	Code    string `json:"code"`    // 登陆返回
	Message string `json:"message"` // 登陆返回错误码
}
