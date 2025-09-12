package esign

import (
	"github.com/renxingcode/esign-go-sdk/api/account_api"
	"github.com/renxingcode/esign-go-sdk/api/auth_api"
	"github.com/renxingcode/esign-go-sdk/api/sign_api"
	"github.com/renxingcode/esign-go-sdk/api/template_api"
	"github.com/renxingcode/esign-go-sdk/config"
)

// e签宝 SDK 的主入口点
type Client struct {
	Auth     *auth_api.AuthService
	Template *template_api.TemplateService
	Account  *account_api.AccountService
	Sign     *sign_api.SignService
}

// NewClient 创建一个新的 e签宝 客户端
func NewClient(cfg *config.Config) *Client {
	// 初始化认证服务
	authService := auth_api.NewAuthService(cfg)
	// 构建并返回主客户端
	return &Client{
		Auth:     authService,
		Template: template_api.NewTemplateService(cfg, authService),
		Account:  account_api.NewAccountService(cfg, authService),
		Sign:     sign_api.NewSignService(cfg, authService),
	}
}
