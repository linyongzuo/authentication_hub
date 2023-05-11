package response

type GenerateCodeResponse struct {
	BaseResp
	Codes []string `json:"codes"`
}
