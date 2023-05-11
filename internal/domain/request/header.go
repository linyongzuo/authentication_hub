package request

type Header struct {
	Version     string      `json:"version"`
	MessageType MessageType `json:"messageType"`
}
