package response

type AdminLoginResponse struct {
	BaseResp
	UserName string `json:"userName"`
}
type UserLoginResponse struct {
	BaseResp
	Mac string `json:"mac"`
}
