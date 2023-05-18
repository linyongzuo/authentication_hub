package request

type MessageType int32

const (
	MessageUnknown           MessageType = 0
	MessageHeartbeat         MessageType = 1
	MessageAdminLogin        MessageType = 2
	MessageAdminGenerateCode MessageType = 3
	MessageUserLogin         MessageType = 4
	MessageUserLogout        MessageType = 5
	MessageUserInfo          MessageType = 6
)
