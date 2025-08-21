package template_api

import (
	"context"
	"github.com/renxingcode/esign-go-sdk/api/auth_api"
	"github.com/renxingcode/esign-go-sdk/config"
	"net/http"
)

// Service 模板服务
type Service struct {
	config      *config.Config
	httpClient  *http.Client
	authService *auth_api.Service // 持有认证服务的引用，用于获取 token
}

// NewTemplateService 创建模板服务实例
func NewTemplateService(cfg *config.Config, authSvc *auth_api.Service) *Service {
	return &Service{
		config:      cfg,
		httpClient:  http.DefaultClient,
		authService: authSvc,
	}
}

// GetTemplateList 获取合同模板列表
func (s *Service) GetTemplateList(ctx context.Context) (string, error) {
	return "test_template_list", nil
}
