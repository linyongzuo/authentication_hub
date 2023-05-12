package request

type AdminLoginReq struct {
	Header
	UserName string `json:"userName"` // 用户名
	Password string `json:"password"` // 密码
}

type UserLoginReq struct {
	Header
	Code string `json:"code"` // 生成的code
}
