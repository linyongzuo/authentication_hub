package request

type HeartbeatReq struct {
	Header
	Admin bool `json:"admin"`
}
