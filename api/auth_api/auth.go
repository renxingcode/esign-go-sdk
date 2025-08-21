package auth_api

import (
	"context"
	"github.com/renxingcode/esign-go-sdk/config"
	"net/http"
	"time"
)

// Service 认证服务, 它应该被嵌入到主 Client 中，并且负责 token 的管理（获取、刷新、缓存）
type Service struct {
	config    *config.Config
	client    *http.Client // 可以接收一个自定义的 http.Client
	token     string
	expiresAt time.Time
}

// NewAuthService 创建一个认证服务实例
func NewAuthService(cfg *config.Config) *Service {
	return &Service{
		config: cfg,
		client: http.DefaultClient, // 默认使用 http.DefaultClient
	}
}

// GetToken 获取有效的 token
func (s *Service) GetToken(ctx context.Context) (string, error) {
	return "test_token", nil
}
