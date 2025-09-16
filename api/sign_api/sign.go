package sign_api

import (
	"errors"
	"fmt"
	"github.com/renxingcode/esign-go-sdk/api"
	"github.com/renxingcode/esign-go-sdk/api/account_api"
	"github.com/renxingcode/esign-go-sdk/api/auth_api"
	"github.com/renxingcode/esign-go-sdk/api/template_api"
	"github.com/renxingcode/esign-go-sdk/config"
	"github.com/renxingcode/esign-go-sdk/types"
	"github.com/renxingcode/esign-go-sdk/utils"
	"net/http"
	"net/url"
	"strings"
)

// SignServiceInterface 流程服务接口
type SignServiceInterface interface {
	ESignCreateFlowOneStep(requestESignCreateFlowData *types.ESignCreateFlowRequestData, writeLog bool) (eSignResponse *types.ESignCommonResponse, err error)
}

var _ SignServiceInterface = (*SignService)(nil)

// SignService 流程服务
type SignService struct {
	config          *config.Config
	httpClient      *http.Client
	authService     *auth_api.AuthService // 持有认证服务的引用，用于获取 token
	templateService *template_api.TemplateService
	accountService  *account_api.AccountService
}

// NewTemplateService 创建模板服务实例
func NewSignService(cfg *config.Config, authSvc *auth_api.AuthService) *SignService {
	return &SignService{
		config:          cfg,
		httpClient:      http.DefaultClient,
		authService:     authSvc,
		templateService: template_api.NewTemplateService(cfg, authSvc),
		accountService:  account_api.NewAccountService(cfg, authSvc),
	}
}

// ESignCreateFlowOneStep 请求e签宝发起签署流程
// e签宝官方接口文档 https://open.esign.cn/doc/opendoc/paas_api/pwd6l4
func (s *SignService) ESignCreateFlowOneStep(requestESignCreateFlowData *types.ESignCreateFlowRequestData, writeLog bool) (eSignResponse *types.ESignCommonResponse, err error) {
	actionName := "一步发起签署:"

	//数据校验
	signerName := requestESignCreateFlowData.SignerName
	signerPhone := requestESignCreateFlowData.SignerPhone
	//companySealID := requestESignCreateFlowData.CompanySealID
	contractFiles := requestESignCreateFlowData.ContractFiles
	if signerName == "" || signerPhone == "" || len(contractFiles) == 0 {
		return nil, errors.New(actionName + "ESignCreateFlowOneStep传入的参数错误:签署人姓名、签署人手机号、合同文件列表都不能为空")
	}

	//Docs
	docsList := make([]types.ESignCreateFlowDocs, 0)
	//Signers
	signersList := make([]types.ESignCreateFlowSigner, 0)
	//公共变量
	thirdOrderNo := "your_own_defined_third_order_no" //改为你自己业务需要的内容
	signerAccountId, err := s.accountService.GetOrCreateESignSignerAccountId(signerName, signerPhone, false)
	if err != nil {
		return nil, err
	}
	if signerAccountId == "" {
		return nil, errors.New("签署人账号获取失败")
	}

	for _, contractFile := range contractFiles {
		fileId := contractFile.EFileId
		if fileId == "" {
			return nil, errors.New("合同文件e_fileid不能为空")
		}
		templateId := contractFile.TemplateId
		if templateId == "" {
			return nil, errors.New("合同模板template_id不能为空")
		}

		//docs
		docsList = append(docsList, types.ESignCreateFlowDocs{
			FileID: fileId,
		})

		// 获取e签宝合同模板信息
		eSignTemplateData, err := s.templateService.GetAndParseESignTemplateDetailData(templateId, true, false)
		if err != nil {
			return nil, err
		}

		var platformSign bool
		signFieldsCompanyList := make([]types.EsignSignField, 0)
		signFieldsUserList := make([]types.EsignSignField, 0)
		for _, Participants := range eSignTemplateData.Participants {
			for _, Components := range Participants.Components {
				if Components.ComponentType == 6 { //签署区
					//公司盖章区域company+_拼接
					componentKeySlice := strings.Split(Components.ComponentKey, "_")
					if len(componentKeySlice) == 0 {
						return nil, errors.New(templateId + ": 合同模板签署区配置异常:应该是以下划线分隔的字符串,例如:company_xxx或者user_xxx")
					}

					signFieldsList := make([]types.EsignSignField, 0)
					if componentKeySlice[0] == "company" { //机构签署
						platformSign = true
						signFieldsItem := types.EsignSignField{
							AutoExecute:        true,
							ActorIndentityType: 2, //机构签约类别，当签约主体为机构时必传：2-机构盖章
							FileID:             fileId,
							//SealID:             requestESignCreateFlowData.CompanySealID, //公司印章ID,如果注释这行,e签宝就回用默认的公司盖章ID处理
							PosBean: types.EsignPosBean{
								PosPage: Components.ComponentPosition.ComponentPageNum,
								PosX:    Components.ComponentPosition.ComponentPositionX,
								PosY:    Components.ComponentPosition.ComponentPositionY,
							},
							SignDateBeanType: 1, //是否需要添加签署日期，0-禁止 1-必须 2-不限制，默认0
						}
						signFieldsList = append(signFieldsList, signFieldsItem)
						signFieldsCompanyList = append(signFieldsCompanyList, signFieldsItem)
					} else { //个人签署
						platformSign = false
						signFieldsItem := types.EsignSignField{
							AutoExecute: false,
							FileID:      fileId,
							PosBean: types.EsignPosBean{
								PosPage: Components.ComponentPosition.ComponentPageNum,
								PosX:    Components.ComponentPosition.ComponentPositionX,
								PosY:    Components.ComponentPosition.ComponentPositionY,
							},
							SignDateBeanType: 2, //是否需要添加签署日期，0-禁止 1-必须 2-不限制，默认0
						}
						signFieldsList = append(signFieldsList, signFieldsItem)
						signFieldsUserList = append(signFieldsUserList, signFieldsItem)
					}

					//types.ESignCreateFlowSigner
					signersItemData := types.ESignCreateFlowSigner{
						PlatformSign: platformSign,
						SignerAccount: types.SignerAccount{
							SignerAccountID: signerAccountId,
						},
						SignFields:   signFieldsList,
						ThirdOrderNo: thirdOrderNo,
					}
					signersList = append(signersList, signersItemData)

				}
			}
		}
		if len(signFieldsCompanyList) == 0 {
			return nil, errors.New(fmt.Sprintf("合同模板缺少company_xxx控件编码;合同模板ID:%s,合同模板名称:%s", templateId, eSignTemplateData.SignTemplateName))
		}
		if len(signFieldsUserList) == 0 {
			return nil, errors.New(fmt.Sprintf("合同模板个人签署区缺少user_xxx控件编码;合同模板ID:%s,合同模板名称:%s", templateId, eSignTemplateData.SignTemplateName))
		}
	}

	//e签宝回调通知地址
	callbackUrl := "http://yourdomain.com/esign/call/back/notice"

	// 构建请求参数
	requestBody := &types.ESignCreateFlowRequest{
		Docs: docsList,
		FlowInfo: types.ESignCreateFlowFlowInfo{
			AutoArchive:   true,
			AutoInitiate:  true,
			BusinessScene: "合同签署", //可以改为你自定义的内容
			FlowConfigInfo: types.EsignFlowConfigInfo{
				NoticeDeveloperUrl: callbackUrl,
			},
		},
		Signers: signersList,
	}

	// 发起HTTP请求
	requestUrl := s.config.BaseURL + api.CreateFlowOneStep
	requestHeaders, err := s.authService.RequestESignHeaders()
	if err != nil {
		return nil, api.BuildRequestESignHeadersError(actionName, err)
	}
	response, err := utils.SendHttpPostRequest(requestUrl, requestBody, requestHeaders, writeLog)
	if err != nil {
		return nil, api.SendHttpRequestError(actionName, err)
	}

	// 解析响应体
	eSignResponse, err = api.GetESignCommonResponse(response)
	if err != nil {
		return nil, api.ParseESignResponseError(actionName, err)
	}

	return eSignResponse, nil
}

// GetESignExecuteUrlByFlowId 查询签署链接
// e签宝官方接口文档 https://open.esign.cn/doc/opendoc/saas_api/fh3gh1_dwz08n
func (s *SignService) GetESignExecuteUrlByFlowId(flowId, signerName, signerPhone string, writeLog bool) (eSignResponse *types.ESignCommonResponse, err error) {
	actionName := "查询签署链接:"

	//数据校验
	if flowId == "" || signerName == "" || signerPhone == "" {
		return nil, errors.New(actionName + "传入的参数错误:flowId、签署人姓名、签署人手机号都不能为空")
	}

	//获取签署人账号
	signerAccountId, err := s.accountService.GetOrCreateESignSignerAccountId(signerName, signerPhone, false)
	if err != nil {
		return nil, err
	}
	if signerAccountId == "" {
		return nil, errors.New("签署人账号获取失败")
	}

	// 构建查询参数
	params := url.Values{}
	params.Add("accountId", signerAccountId)

	// 发起HTTP请求
	//将 api.GetESignExecuteUrlByFlowId 中的 {FLOW_ID} 替换为 flowId
	requestPath := strings.Replace(api.GetESignExecuteUrlByFlowId, "{FLOW_ID}", flowId, 1)
	requestUrl := s.config.BaseURL + requestPath + "?" + params.Encode()
	requestHeaders, err := s.authService.RequestESignHeaders()
	if err != nil {
		return nil, api.BuildRequestESignHeadersError(actionName, err)
	}
	response, err := utils.SendHttpGetRequest(requestUrl, requestHeaders, writeLog)
	if err != nil {
		return nil, api.SendHttpRequestError(actionName, err)
	}

	// 解析响应体
	eSignResponse, err = api.GetESignCommonResponse(response)
	if err != nil {
		return nil, api.ParseESignResponseError(actionName, err)
	}
	return eSignResponse, nil
}

// ESignFlowRevoke PUT 撤回签署流程:撤销签署流程，撤销后流程中止，所有签署短信打开后无效。
// 文档地址: https://open.esign.cn/doc/opendoc/saas_api/hv1dii_uqoamg
func (s *SignService) ESignFlowRevoke(flowId string, writeLog bool) (eSignResponse *types.ESignCommonResponse, err error) {
	actionName := "撤回签署流程:"

	//数据校验
	if flowId == "" {
		return nil, errors.New(actionName + "传入的参数错误:flowId不能为空")
	}

	// 发起HTTP请求,这里需要PUT请求
	requestPath := strings.Replace(api.ESignFlowRevokePath, "{FLOW_ID}", flowId, 1) //替换 {FLOW_ID}
	requestUrl := s.config.BaseURL + requestPath
	requestHeaders, err := s.authService.RequestESignHeaders()
	if err != nil {
		return nil, api.BuildRequestESignHeadersError(actionName, err)
	}
	response, err := utils.SendHttpPutRequest(requestUrl, requestHeaders, writeLog)
	if err != nil {
		return nil, api.SendHttpRequestError(actionName, err)
	}

	// 解析响应体
	eSignResponse, err = api.GetESignCommonResponse(response)
	if err != nil {
		return nil, api.ParseESignResponseError(actionName, err)
	}
	return eSignResponse, nil
}
