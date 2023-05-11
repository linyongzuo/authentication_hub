package request

type GenerateCodeReq struct {
	Header
	Count int `json:"count"` // 生成个数
}
