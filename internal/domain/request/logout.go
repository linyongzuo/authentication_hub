package request

type UserLogoutReq struct {
	Header
	Mac    string `json:"mac"` // 如果是接收到用户下线消息
	Ip     string `json:"ip"`  // 如果是接收到系统下线消息
	System bool   `json:"system"`
}
