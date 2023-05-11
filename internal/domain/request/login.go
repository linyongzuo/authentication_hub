package request

type LoginReq struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}
