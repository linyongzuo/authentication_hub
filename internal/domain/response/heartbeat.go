package response

type HeartbeatResp struct {
	BaseResp
	Mac string `json:"mac"`
}
