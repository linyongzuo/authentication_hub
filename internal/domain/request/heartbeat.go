package request

type HeartbeatReq struct {
	Header
	Mac string `json:"mac"`
	Ip  string `json:"ip"`
}
