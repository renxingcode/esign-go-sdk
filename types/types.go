package types

// ESignCommonResponse e签宝通用响应结构体
type ESignCommonResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// GetESignTokenRequest 获取e签宝token的请求体
type GetESignTokenRequest struct {
	AppId     string `json:"appId"`
	Secret    string `json:"secret"`
	GrantType string `json:"grantType,default:client_credentials"`
}

type GetESignTokenResponse struct {
	ExpiresIn    string `json:"expiresIn"`
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}
