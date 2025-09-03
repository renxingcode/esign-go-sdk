package auth_api

import (
	"errors"
	"fmt"
	"github.com/renxingcode/esign-go-sdk/api"
	"github.com/renxingcode/esign-go-sdk/config"
	"github.com/renxingcode/esign-go-sdk/types"
	"github.com/renxingcode/esign-go-sdk/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"strconv"
	"time"
)

// AuthServiceInterface 认证服务应实现的接口
type AuthServiceInterface interface {
	GetESignToken(useCache bool) (token string, err error)
}

// 编译时检查 Service 是否实现了 AuthServiceInterface 接口
var _ AuthServiceInterface = (*AuthService)(nil)

// AuthService 认证服务, 它应该被嵌入到主 Client 中，并且负责 token 的管理（获取、刷新、缓存）
type AuthService struct {
	config    *config.Config
	client    *http.Client // 可以接收一个自定义的 http.Client
	expiresAt time.Time
	token     string
}

// NewAuthService 创建一个认证服务实例
func NewAuthService(cfg *config.Config) *AuthService {
	return &AuthService{
		config: cfg,
		client: http.DefaultClient, // 默认使用 http.DefaultClient
	}
}

// GetESignToken 获取e签宝的token
func (s *AuthService) GetESignToken(useCache bool) (token string, err error) {
	if useCache {
		// 从缓存中获取token
		token, err = s.GetESignTokenFromCacheData()
		if err != nil {
			return "", err
		}
		if token != "" {
			logx.Infof("从缓存获取token成功")
			return token, nil
		}
	}

	// 从e签宝服务器获取token
	eSignTokenResp, err := s.GetESignTokenFromESignServer()
	if err != nil {
		return "", err
	}
	token = eSignTokenResp.Token
	logx.Infof("从e签宝服务器获取token成功")

	// 缓存token
	err = s.SetESignTokenToCacheData(token, eSignTokenResp.ExpiresIn)
	if err != nil {
		return token, err
	}
	logx.Infof("写入缓存token成功")

	return token, nil
}

func (s *AuthService) GetESignTokenFromESignServer() (eSignTokenResp *types.GetESignTokenResponse, err error) {
	// 发起HTTP请求,获取e签宝的token
	requestUrl := s.config.BaseURL + api.GetESignTokenPath
	requestBody := &types.GetESignTokenRequest{
		AppId:     s.config.AppID,
		Secret:    s.config.AppSecret,
		GrantType: s.config.GrantType,
	}
	requestHeaders := map[string]string{
		"Content-Type": "application/json; charset=UTF-8",
	}
	response, err := utils.SendHttpPostRequest(requestUrl, requestBody, requestHeaders, s.config.IsWriteLog)
	if err != nil {
		return nil, err
	}

	// 解析响应体
	responseStruct, err := api.GetESignCommonResponse(response)
	if err != nil {
		return nil, err
	}

	// 解析Data结构
	responseDataStruct := types.GetESignTokenResponse{}
	err = utils.JsonUnmarshalToStruct(responseStruct.Data, &responseDataStruct)
	if err != nil {
		return nil, err
	}
	return &responseDataStruct, nil
}

func (s *AuthService) GetESignTokenFromCacheData() (token string, err error) {
	if s.config.RedisClient == nil {
		logx.Infof("Redis组件未初始化，跳过缓存获取")
		return "", nil
	}

	token, err = s.config.RedisClient.Get(api.ESignAccessTokenKey)
	if err != nil {
		logx.Errorf("从缓存获取 token 失败: %v", err)
		return "", err
	}

	//调用e签宝接口检测Token是否有效,如果无效则重新获取 todo

	return token, nil
}

func (s *AuthService) SetESignTokenToCacheData(token string, eSignExpiresIn string) error {
	if s.config.RedisClient == nil {
		logx.Infof("Redis组件未初始化，跳过缓存设置")
		return nil
	}

	//计算过期时间
	//expiresIn 是毫秒级时间戳，需用 expiresIn - 当前时间 得到剩余有效期
	expiresInMs, err := strconv.ParseInt(eSignExpiresIn, 10, 64)
	if err != nil {
		return fmt.Errorf("解析e签宝的过期时间失败: %w", err)
	}
	nowMs := time.Now().UnixMilli()
	remainMs := expiresInMs - nowMs - 60*1000 // 提前1分钟失效
	if remainMs <= 0 {
		remainMs = 5 * 60 * 1000 // 兜底5分钟
	}
	expireDuration := time.Duration(remainMs) * time.Millisecond
	logx.Infof("设置缓存token的过期时间: %v", expireDuration)

	//设置缓存
	err = s.config.RedisClient.Set(api.ESignAccessTokenKey, token, expireDuration)
	if err != nil {
		logx.Errorf("设置缓存token失败: %v", err)
		return err
	}
	return nil
}

func (s *AuthService) RequestESignHeaders() (map[string]string, error) {
	token, err := s.GetESignToken(false)
	if err != nil {
		return nil, errors.New("获取e签宝token失败:" + err.Error())
	}
	requestHeaders := map[string]string{
		"Content-Type":        "application/json; charset=UTF-8",
		"X-Tsign-Open-App-Id": s.config.AppID,
		"X-Tsign-Open-Token":  token,
	}
	return requestHeaders, nil
}
