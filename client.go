package esign

import (
	"github.com/renxingcode/esign-go-sdk/api/auth_api"
	"github.com/renxingcode/esign-go-sdk/api/template_api"
	"github.com/renxingcode/esign-go-sdk/config"
)

// e签宝 SDK 的主入口点
type Client struct {
	Auth     *auth_api.Service
	Template *template_api.Service
}

// NewClient 创建一个新的 e签宝 客户端
func NewClient(cfg *config.Config) *Client {
	// 初始化认证服务
	authService := auth_api.NewAuthService(cfg)

	// 初始化模板服务，并注入依赖（配置和认证服务）
	templateService := template_api.NewTemplateService(cfg, authService)

	// 构建并返回主客户端
	return &Client{
		Auth:     authService,
		Template: templateService,
	}
}
