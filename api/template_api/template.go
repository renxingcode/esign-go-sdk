package template_api

import (
	"context"
	"github.com/renxingcode/esign-go-sdk/api/auth_api"
	"github.com/renxingcode/esign-go-sdk/config"
	"net/http"
)

// TemplateServiceInterface 模板服务接口
type TemplateServiceInterface interface {
	GetESignTemplateList(ctx context.Context) (templateList string, err error)
}

var _ TemplateServiceInterface = (*TemplateService)(nil)

// TemplateService 模板服务
type TemplateService struct {
	config      *config.Config
	httpClient  *http.Client
	authService *auth_api.AuthService // 持有认证服务的引用，用于获取 token
}

// NewTemplateService 创建模板服务实例
func NewTemplateService(cfg *config.Config, authSvc *auth_api.AuthService) *TemplateService {
	return &TemplateService{
		config:      cfg,
		httpClient:  http.DefaultClient,
		authService: authSvc,
	}
}

// GetESignTemplateList 获取合同模板列表
func (s *TemplateService) GetESignTemplateList(ctx context.Context) (templateList string, err error) {
	return "test_template_list", nil
}
