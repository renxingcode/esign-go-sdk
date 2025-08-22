package auth_api

import (
	"context"
	"github.com/renxingcode/esign-go-sdk/api"
	"github.com/renxingcode/esign-go-sdk/config"
	"github.com/renxingcode/esign-go-sdk/types"
	"github.com/renxingcode/esign-go-sdk/utils"
	"net/http"
	"time"
)

// AuthServiceInterface 认证服务应实现的接口
type AuthServiceInterface interface {
	GetESignToken(ctx context.Context) (token string, err error)
}

// 编译时检查 Service 是否实现了 AuthServiceInterface 接口
var _ AuthServiceInterface = (*AuthService)(nil)

// AuthService 认证服务, 它应该被嵌入到主 Client 中，并且负责 token 的管理（获取、刷新、缓存）
type AuthService struct {
	config    *config.Config
	client    *http.Client // 可以接收一个自定义的 http.Client
	token     string
	expiresAt time.Time
}

// NewAuthService 创建一个认证服务实例
func NewAuthService(cfg *config.Config) *AuthService {
	return &AuthService{
		config: cfg,
		client: http.DefaultClient, // 默认使用 http.DefaultClient
	}
}

// GetESignToken 获取e签宝的token
func (s *AuthService) GetESignToken(ctx context.Context) (token string, err error) {
	// 发起HTTp请求,获取e签宝的token
	requestUrl := s.config.BaseURL + api.GetESignTokenPath
	requestBody := &types.GetESignTokenRequest{
		AppId:     s.config.AppID,
		Secret:    s.config.AppSecret,
		GrantType: s.config.GrantType,
	}
	response, err := utils.SendHttpPostRequest(requestUrl, requestBody, nil, s.config.IsWriteLog)
	if err != nil {
		return "", err
	}

	// 解析响应体
	responseStruct, err := api.GetESignCommonResponse(response)
	if err != nil {
		return "", err
	}

	// 解析Data结构
	responseDataStruct := &types.GetESignTokenResponse{}
	err = utils.JsonUnmarshalToStruct(responseStruct.Data, &responseDataStruct)
	if err != nil {
		return "", err
	}
	return responseDataStruct.Token, nil
}
