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
	companySealID := requestESignCreateFlowData.CompanySealID
	contractFiles := requestESignCreateFlowData.ContractFiles
	if signerName == "" || signerPhone == "" || companySealID == "" || len(contractFiles) == 0 {
		return nil, errors.New(actionName + "ESignCreateFlowOneStep传入的参数错误:签署人姓名、签署人手机号、公司印章ID、合同文件列表都不能为空")
	}

	//Docs
	docsList := make([]types.ESignCreateFlowDocs, 0)
	//Signers
	signersList := make([]types.ESignCreateFlowSigner, 0)
	//公共变量
	thirdOrderNo := "your_own_defined_third_order_no" //改为你自己业务需要的内容
	signerAccountId, err := s.accountService.GetOrCreateESignSignerAccountId(signerName, signerPhone, writeLog)
	if err != nil {
		return nil, err
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
				if Components.ComponentType == 6 {
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
							ActorIndentityType: 2,
							FileID:             fileId,
							SealID:             companySealID, //公司印章ID,如果产生阻塞,注释这行就可以了,e签宝就回用默认的公司盖章ID处理
							PosBean: types.EsignPosBean{
								PosPage: Components.ComponentPosition.ComponentPageNum,
								PosX:    Components.ComponentPosition.ComponentPositionX,
								PosY:    Components.ComponentPosition.ComponentPositionY,
							},
							SignDateBeanType: 1,
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
							SignDateBeanType: 2,
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
	eSignResponse, err = api.GetESignCommonResponse(response)
	if err != nil {
		return nil, api.ParseESignResponseError(actionName, err)
	}

	return eSignResponse, nil
}
