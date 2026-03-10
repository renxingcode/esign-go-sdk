package account_api

import (
	"github.com/renxingcode/esign-go-sdk/api"
	"github.com/renxingcode/esign-go-sdk/api/auth_api"
	"github.com/renxingcode/esign-go-sdk/config"
	"github.com/renxingcode/esign-go-sdk/types"
	"github.com/renxingcode/esign-go-sdk/utils"
	"net/http"
	"net/url"
	"strings"
)

// AccountServiceInterface 账户服务接口
type AccountServiceInterface interface {
	GetESignPersonsIdentityInfo(psnAccount string, writeLog bool) (eSignPersonsIdentityData *types.PersonsIdentityData, err error)
	CreateESignPersonsIdentity(name, mobile, thirdPartyUserId string, writeLog bool) (eSignPersonsIdentityData *types.CreateESignPersonsIdentityResponse, err error)
	GetOrCreateESignSignerAccountId(name, mobile string, writeLog bool) (accountId string, err error)
}

var _ AccountServiceInterface = (*AccountService)(nil)

// AccountService 账户服务
type AccountService struct {
	config      *config.Config
	httpClient  *http.Client
	authService *auth_api.AuthService // 持有认证服务的引用，用于获取 token
}

// NewTemplateService 创建模板服务实例
func NewAccountService(cfg *config.Config, authSvc *auth_api.AuthService) *AccountService {
	return &AccountService{
		config:      cfg,
		httpClient:  http.DefaultClient,
		authService: authSvc,
	}
}

// GetESignPersonsIdentityInfo 查询个人认证信息
// e签宝官方接口文档 https://open.esign.cn/doc/opendoc/auth3/vssvtu
func (s *AccountService) GetESignPersonsIdentityInfo(psnAccount string, writeLog bool) (eSignPersonsIdentityData *types.PersonsIdentityData, err error) {
	actionName := "查询个人认证信息:"

	// 构建查询参数
	params := url.Values{}
	params.Set("psnAccount", psnAccount)

	// 发起HTTP请求
	requestUrl := s.config.BaseURL + api.GetPersonsIdentityInfo + "?" + params.Encode()
	requestHeaders, err := s.authService.RequestESignHeaders()
	if err != nil {
		return nil, api.BuildRequestESignHeadersError(actionName, err)
	}
	response, err := utils.SendHttpGetRequest(requestUrl, requestHeaders, writeLog)
	if err != nil {
		return nil, api.SendHttpRequestError(actionName, err)
	}

	// 解析响应体
	eSignResponse, err := api.GetESignCommonResponse(response)
	if err != nil {
		return nil, api.ParseESignResponseError(actionName, err)
	}
	if eSignResponse.Code != api.ESignResponseCodeSuccess {
		return nil, api.GetESignResponseError(eSignResponse)
	}

	// 解析Data结构
	err = utils.JsonUnmarshalToStruct(eSignResponse.Data, &eSignPersonsIdentityData)
	if err != nil {
		return nil, api.ParseESignResponseDataError(actionName, err)
	}
	return eSignPersonsIdentityData, nil
}

// CreateESignPersonsIdentity 创建个人认证信息
// e签宝官方接口文档 https://open.esign.cn/doc/opendoc/paas_api/ox6nog
func (s *AccountService) CreateESignPersonsIdentity(name, mobile, thirdPartyUserId string, writeLog bool) (eSignPersonsIdentityData *types.CreateESignPersonsIdentityResponse, err error) {
	actionName := "创建个人认证信息:"

	// 构建请求参数
	requestBody := &types.CreateESignPersonsIdentityRequest{
		Name:             name,
		Mobile:           mobile,
		ThirdPartyUserId: thirdPartyUserId,
	}

	// 发起HTTP请求
	requestUrl := s.config.BaseURL + api.CreatePersonsIdentity
	requestHeaders, err := s.authService.RequestESignHeaders()
	if err != nil {
		return nil, api.BuildRequestESignHeadersError(actionName, err)
	}
	response, err := utils.SendHttpPostRequest(requestUrl, requestBody, requestHeaders, writeLog)
	if err != nil {
		return nil, api.SendHttpRequestError(actionName, err)
	}

	// 解析响应体
	eSignResponse, err := api.GetESignCommonResponse(response)
	if err != nil {
		return nil, api.ParseESignResponseError(actionName, err)
	}

	// 个人认证信息不存在，创建个人认证信息并返回;如果存在,忽略错误提示,直接返回accountId
	//if eSignResponse.Code != api.ESignResponseCodeSuccess {
	//	return nil, api.GetESignResponseError(eSignResponse)
	//}
	if eSignResponse.Data == nil {
		return nil, api.GetESignResponseError(eSignResponse)
	}

	// 解析Data结构
	err = utils.JsonUnmarshalToStruct(eSignResponse.Data, &eSignPersonsIdentityData)
	if err != nil {
		return nil, api.ParseESignResponseDataError(actionName, err)
	}
	return eSignPersonsIdentityData, nil
}

// UpdateESignPersonsIdentity 修改个人认证信息
// e签宝官方接口文档 https://open.esign.cn/doc/opendoc/saas_api/tzi4kd_ma1d8m
func (s *AccountService) UpdateESignPersonsIdentity(signerAccountId string, requestData map[string]string, writeLog bool) (eSignPersonsIdentityData *types.UpdateESignPersonsIdentityResponse, err error) {
	actionName := "修改个人认证信息:"

	// 发起HTTP请求
	requestPath := strings.Replace(api.UpdatePersonsIdentity, "{ACCOUNT_ID}", signerAccountId, 1)
	requestUrl := s.config.BaseURL + requestPath
	requestHeaders, err := s.authService.RequestESignHeaders()
	if err != nil {
		return nil, api.BuildRequestESignHeadersError(actionName, err)
	}
	response, err := utils.SendHttpPutRequest(requestUrl, requestData, requestHeaders, writeLog)
	if err != nil {
		return nil, api.SendHttpRequestError(actionName, err)
	}

	// 解析响应体
	eSignResponse, err := api.GetESignCommonResponse(response)
	if err != nil {
		return nil, api.ParseESignResponseError(actionName, err)
	}

	if eSignResponse.Data == nil {
		return nil, api.GetESignResponseError(eSignResponse)
	}

	// 解析Data结构
	err = utils.JsonUnmarshalToStruct(eSignResponse.Data, &eSignPersonsIdentityData)
	if err != nil {
		return nil, api.ParseESignResponseDataError(actionName, err)
	}
	return eSignPersonsIdentityData, nil
}

// GetESignPersonsIdentityInfo 查询个人认证信息
// e签宝官方接口文档 https://open.esign.cn/doc/opendoc/auth3/vssvtu
func (s *AccountService) GetOrCreateESignSignerAccountId(name, mobile string, writeLog bool) (accountId string, err error) {
	psnAccount := mobile

	//根据手机号获取accountId 直接通过 /v1/accounts/createByThirdPartyUserId 就可以了，
	//如果不存在自动创建，如果存在则忽略code：53000000，直接拿返回的 data.accountId 用就可以了。
	//因此,不需要调用 /v3/persons/identity-info 这个接口去获取accountId.
	/*
		// 先查询个人认证信息
		getESignPersonsIdentityData, err := s.GetESignPersonsIdentityInfo(psnAccount, writeLog)
		if err != nil {
			return "", err
		}
		if getESignPersonsIdentityData.PsnId != "" {
			return getESignPersonsIdentityData.PsnId, nil
		}
	*/

	// 个人认证信息不存在，创建个人认证信息并返回;如果存在,忽略错误提示,直接返回accountId
	createESignPersonsIdentityData, err := s.CreateESignPersonsIdentity(name, mobile, psnAccount, writeLog)
	//fmt.Println("createESignPersonsIdentityData:", utils.JsonMarshalIndent(createESignPersonsIdentityData))
	if createESignPersonsIdentityData != nil && createESignPersonsIdentityData.AccountId != "" {
		return createESignPersonsIdentityData.AccountId, nil
	}
	return "", nil
}
