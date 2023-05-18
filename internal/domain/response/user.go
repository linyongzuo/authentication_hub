package response

type UserInfoResponse struct {
	BaseResp
	RegisterCount int `json:"registerCount"`
	OnlineCount   int `json:"onlineCount"`
	OfflineCount  int `json:"offlineCount"`
}
