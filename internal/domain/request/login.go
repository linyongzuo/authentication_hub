package request

type AdminLoginReq struct {
	Header
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type UserLoginReq struct {
	Header
	Mac  string `json:"mac"`
	Ip   string `json:"ip"`
	Code string `json:"code"`
}
