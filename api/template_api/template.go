package template_api

import (
	"errors"
	"github.com/renxingcode/esign-go-sdk/api"
	"github.com/renxingcode/esign-go-sdk/api/auth_api"
	"github.com/renxingcode/esign-go-sdk/config"
	"github.com/renxingcode/esign-go-sdk/types"
	"github.com/renxingcode/esign-go-sdk/utils"
	"net/http"
	"net/url"
	"strconv"
)

// TemplateServiceInterface 模板服务接口
type TemplateServiceInterface interface {
	GetESignTemplateDetail(eSignTemplateId string, queryComponents bool, writeLog bool) (eSignResponse *types.ESignCommonResponse, err error)
	CreateByTemplate(eSignTemplateId, name string, simpleFormFields map[string]string, writeLog bool) (eSignResponse *types.ESignCommonResponse, err error)
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

// GetESignTemplateDetail 获取流程模版详情
// e签宝官方接口文档 https://open.esign.cn/doc/opendoc/file-and-template3/pfzut7ho9obc7c5r
func (s *TemplateService) GetESignTemplateDetail(eSignTemplateId string, queryComponents bool, writeLog bool) (eSignResponse *types.ESignCommonResponse, err error) {
	actionName := "请求e签宝获取流程模版详情"

	// 构建查询参数
	params := url.Values{}
	params.Add("orgId", s.config.OrgId)
	params.Add("signTemplateId", eSignTemplateId)
	if qc := strconv.FormatBool(queryComponents); qc != "" { //布尔值转字符串
		params.Add("queryComponents", qc)
	}

	// 发起HTTP请求
	requestUrl := s.config.BaseURL + api.GetESignTemplateDetailPath + "?" + params.Encode()
	requestHeaders, err := s.authService.RequestESignHeaders()
	if err != nil {
		return nil, errors.New(actionName + "构建请求e签宝的headers失败:" + err.Error())
	}
	response, err := utils.SendHttpGetRequest(requestUrl, requestHeaders, writeLog)
	if err != nil {
		return nil, errors.New(actionName + "发送http请求失败:" + err.Error())
	}

	// 解析响应体
	eSignResponse, err = api.GetESignCommonResponse(response)
	if err != nil {
		return nil, errors.New(actionName + "解析e签宝响应体失败:" + err.Error())
	}
	return eSignResponse, nil
}

// CreateByTemplate 通过模板创建文件
// e签宝官方接口文档 https://open.esign.cn/doc/opendoc/saas_api/cz9d65_sh823i?searchText=
func (s *TemplateService) CreateByTemplate(eSignTemplateDocFileId, eSignTemplateDocFileName string, simpleFormFields map[string]string, writeLog bool) (eSignResponse *types.ESignCommonResponse, err error) {
	actionName := "通过模板创建文件"

	requestUrl := s.config.BaseURL + api.CreateESignFileByTemplate
	requestBody := &types.CreateESignFileByTemplateRequest{
		TemplateDocFileId:   eSignTemplateDocFileId,
		TemplateDocFileName: eSignTemplateDocFileName,
		SimpleFormFields:    simpleFormFields,
	}
	requestHeaders, err := s.authService.RequestESignHeaders()
	if err != nil {
		return nil, errors.New(actionName + "构建请求e签宝的headers失败:" + err.Error())
	}
	response, err := utils.SendHttpPostRequest(requestUrl, requestBody, requestHeaders, writeLog)
	if err != nil {
		return nil, errors.New(actionName + "发送http请求失败:" + err.Error())
	}

	// 解析响应体
	eSignResponse, err = api.GetESignCommonResponse(response)
	if err != nil {
		return nil, errors.New(actionName + "解析e签宝响应体失败:" + err.Error())
	}
	return eSignResponse, nil
}
