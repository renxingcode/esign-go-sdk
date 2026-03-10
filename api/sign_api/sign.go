package sign_api

import (
	"context"
	"errors"
	"fmt"
	"github.com/renxingcode/esign-go-sdk/api"
	"github.com/renxingcode/esign-go-sdk/api/account_api"
	"github.com/renxingcode/esign-go-sdk/api/auth_api"
	"github.com/renxingcode/esign-go-sdk/api/template_api"
	"github.com/renxingcode/esign-go-sdk/config"
	"github.com/renxingcode/esign-go-sdk/types"
	"github.com/renxingcode/esign-go-sdk/utils"
	"golang.org/x/sync/errgroup"
	"net/http"
	"net/url"
	"strings"
)

// SignServiceInterface 流程服务接口
type SignServiceInterface interface {
	ESignCreateFlowOneStep(requestESignCreateFlowData *types.ESignCreateFlowRequestData, writeLog bool) (eSignResponse *types.ESignCommonResponse, err error)
	GetESignExecuteUrlByFlowId(flowId, signerName, signerPhone string, writeLog bool) (eSignResponse *types.ESignCommonResponse, err error)
	ESignFlowRevoke(flowId string, writeLog bool) (eSignResponse *types.ESignCommonResponse, err error)
	GetESignDocumentsUrlByFlowId(flowId string, writeLog bool) (eSignDocumentsDocs []types.GetDocumentsUrlResponseDataDocs, err error)
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

// ESignCreateFlowOneStep 请求e签宝发起签署流程,此流程为e签宝自动给签署人发送短信
// e签宝官方接口文档 https://open.esign.cn/doc/opendoc/paas_api/pwd6l4
func (s *SignService) ESignCreateFlowOneStep(requestData *types.ESignCreateFlowRequestData, writeLog bool) (*types.ESignCommonResponse, error) {
	actionName := "一步发起签署:"

	// ----- 1. 参数校验 -----
	if requestData.SignerName == "" || requestData.SignerPhone == "" || len(requestData.ContractFiles) == 0 {
		return nil, errors.New(actionName + "签署人姓名、手机号、合同文件列表都不能为空")
	}

	// ----- 2. 获取签署人账号ID（串行执行）-----
	signerAccountID, err := s.accountService.GetOrCreateESignSignerAccountId(
		requestData.SignerName, requestData.SignerPhone, false,
	)
	if err != nil {
		return nil, fmt.Errorf("%s获取签署人账号失败: %w", actionName, err)
	}
	if signerAccountID == "" {
		return nil, errors.New(actionName + "签署人账号ID为空")
	}

	// ----- 3. 准备并发处理合同文件 -----
	const maxConcurrent = 10 // 最大并发数，可根据需要调整
	sem := make(chan struct{}, maxConcurrent)

	// 定义结果结构
	type fileResult struct {
		FileID  string
		Signers []types.ESignCreateFlowSigner
	}
	resultCh := make(chan fileResult, len(requestData.ContractFiles))

	g, ctx := errgroup.WithContext(context.Background())
	thirdOrderNo := "your_own_defined_third_order_no" // 业务自定义订单号

	// 为每个合同文件启动一个 goroutine
	for _, cf := range requestData.ContractFiles {
		cf := cf // 捕获循环变量
		g.Go(func() error {
			// 并发控制：获取信号量
			select {
			case sem <- struct{}{}:
				defer func() { <-sem }()
			case <-ctx.Done():
				return ctx.Err()
			}

			// 执行单个文件处理
			fileID, signers, err := s.processSingleFile(ctx, cf, signerAccountID, thirdOrderNo, writeLog)
			if err != nil {
				return err // errgroup 会取消其他任务
			}

			// 发送结果（非阻塞，channel 有缓冲）
			select {
			case resultCh <- fileResult{FileID: fileID, Signers: signers}:
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		})
	}

	// 等待所有 goroutine 完成并关闭 resultCh
	go func() {
		_ = g.Wait() // 等待所有任务完成或出错
		close(resultCh)
	}()

	// ----- 4. 收集并发处理结果 -----
	docsList := make([]types.ESignCreateFlowDocs, 0, len(requestData.ContractFiles))
	signersList := make([]types.ESignCreateFlowSigner, 0)

	for res := range resultCh {
		docsList = append(docsList, types.ESignCreateFlowDocs{FileID: res.FileID})
		signersList = append(signersList, res.Signers...)
	}

	// 检查是否有错误发生
	if err = g.Wait(); err != nil {
		return nil, fmt.Errorf("%s并发处理合同文件失败: %w", actionName, err)
	}

	// ----- 5. 构建最终请求并调用e签宝接口 -----
	callbackURL := "http://yourdomain.com/esign/call/back/notice"
	requestBody := &types.ESignCreateFlowRequest{
		Docs: docsList,
		FlowInfo: types.ESignCreateFlowFlowInfo{
			AutoArchive:   true,
			AutoInitiate:  true,
			BusinessScene: "合同签署",
			FlowConfigInfo: types.EsignFlowConfigInfo{
				NoticeDeveloperUrl: callbackURL,
			},
		},
		Signers: signersList,
	}

	requestURL := s.config.BaseURL + api.CreateFlowOneStep
	headers, err := s.authService.RequestESignHeaders()
	if err != nil {
		return nil, api.BuildRequestESignHeadersError(actionName, err)
	}
	resp, err := utils.SendHttpPostRequest(requestURL, requestBody, headers, writeLog)
	if err != nil {
		return nil, api.SendHttpRequestError(actionName, err)
	}

	eSignResponse, err := api.GetESignCommonResponse(resp)
	if err != nil {
		return nil, api.ParseESignResponseError(actionName, err)
	}
	return eSignResponse, nil
}

// processSingleFile 处理单个合同文件：获取模板详情、解析签署区、生成签署人条目
func (s *SignService) processSingleFile(
	ctx context.Context,
	file types.ESignCreateFlowFiles,
	signerAccountID, thirdOrderNo string,
	writeLog bool,
) (fileID string, signers []types.ESignCreateFlowSigner, err error) {
	// 校验文件必要字段
	if file.EFileId == "" {
		return "", nil, errors.New("合同文件 e_fileid 不能为空")
	}
	if file.TemplateId == "" {
		return "", nil, errors.New("合同模板 template_id 不能为空")
	}

	// 获取模板详情
	templateData, err := s.templateService.GetAndParseESignTemplateDetailData(file.TemplateId, true, false)
	if err != nil {
		return "", nil, fmt.Errorf("获取模板详情失败: %w", err)
	}

	var (
		fileSigners []types.ESignCreateFlowSigner
		hasCompany  bool
		hasUser     bool
	)

	// 遍历模板参与方与组件，解析签署区
	for _, participant := range templateData.Participants {
		for _, comp := range participant.Components {
			if comp.ComponentType != 6 { // 6=签署区
				continue
			}

			// 解析控件编码，判断签署主体类型
			keyParts := strings.Split(comp.ComponentKey, "_")
			if len(keyParts) == 0 {
				return "", nil, fmt.Errorf("模板 %s 签署区控件编码格式错误: %s", file.TemplateId, comp.ComponentKey)
			}

			var platformSign bool
			signField := types.EsignSignField{
				FileID: file.EFileId,
				PosBean: types.EsignPosBean{
					PosPage: comp.ComponentPosition.ComponentPageNum,
					PosX:    comp.ComponentPosition.ComponentPositionX,
					PosY:    comp.ComponentPosition.ComponentPositionY,
				},
				SignDateBeanType: 2, // 默认2（不限制），机构签署时会覆盖
			}

			if keyParts[0] == "company" { // 机构签署
				platformSign = true
				signField.AutoExecute = true
				signField.ActorIndentityType = 2 // 机构盖章
				// signField.SealID = requestData.CompanySealID // 若需指定印章可取消注释
				signField.SignDateBeanType = 1 // 必须包含签署日期
				hasCompany = true
			} else { // 个人签署
				platformSign = false
				signField.AutoExecute = false
				signField.SignDateBeanType = 2 // 不限制
				hasUser = true
			}

			// 每个签署区生成一个独立的签署人条目（与原逻辑一致）
			signerItem := types.ESignCreateFlowSigner{
				PlatformSign: platformSign,
				SignerAccount: types.SignerAccount{
					SignerAccountID: signerAccountID,
				},
				SignFields:   []types.EsignSignField{signField},
				ThirdOrderNo: thirdOrderNo,
			}
			fileSigners = append(fileSigners, signerItem)
		}
	}

	// 校验必须同时包含机构和个人签署区
	if !hasCompany {
		return "", nil, fmt.Errorf("合同模板缺少 company_xxx 控件编码; 模板ID: %s, 模板名称: %s",
			file.TemplateId, templateData.SignTemplateName)
	}
	if !hasUser {
		return "", nil, fmt.Errorf("合同模板个人签署区缺少 user_xxx 控件编码; 模板ID: %s, 模板名称: %s",
			file.TemplateId, templateData.SignTemplateName)
	}

	return file.EFileId, fileSigners, nil
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
		return nil, errors.New(actionName + "签署人账号获取失败")
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
	response, err := utils.SendHttpPutRequest(requestUrl, nil, requestHeaders, writeLog)
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

// GetESignDocumentsUrlByFlowId 获取签署完成后的文档链接
// 文档地址: https://open.esign.cn/doc/opendoc/saas_api/oyqsoq_zknh6g
func (s *SignService) GetESignDocumentsUrlByFlowId(flowId string, writeLog bool) (eSignDocumentsDocs []types.GetDocumentsUrlResponseDataDocs, err error) {
	actionName := "查询签署完成后的文档链接:"

	//数据校验
	if flowId == "" {
		return nil, errors.New(actionName + "传入的参数错误:flowId不能为空")
	}

	// 发起HTTP请求
	requestPath := strings.Replace(api.GetESignDocumentsUrlByFlowId, "{FLOW_ID}", flowId, 1) //替换 {FLOW_ID}
	requestUrl := s.config.BaseURL + requestPath
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
	documentsUrlResponseData := &types.GetDocumentsUrlResponseData{}
	err = utils.JsonUnmarshalToStruct(eSignResponse.Data, &documentsUrlResponseData)
	if err != nil {
		return nil, api.ParseESignResponseDataError(actionName, err)
	}
	return documentsUrlResponseData.Docs, nil
}
