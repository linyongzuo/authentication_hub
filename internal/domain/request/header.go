package request

type Header struct {
	Version     string      `json:"version"`
	MessageType MessageType `json:"messageType"`
	Mac         string      `json:"mac"`
	Ip          string      `json:"ip"`
}
