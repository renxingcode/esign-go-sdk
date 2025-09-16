package tests

import (
	"github.com/renxingcode/esign-go-sdk"
	"github.com/renxingcode/esign-go-sdk/Initialize"
	"github.com/renxingcode/esign-go-sdk/types"
	"github.com/renxingcode/esign-go-sdk/utils"
	"testing"
)

var signerName = "张三"
var signerPhone = "13212345678"

// TestESignCreateFlowOneStep 请求e签宝发起签署流程 | go test tests/sign_test.go -v -run TestESignCreateFlowOneStep
func TestESignCreateFlowOneStep(t *testing.T) {
	testClient, err := Initialize.NewTestClient()
	if err != nil {
		t.Errorf("创建测试客户端失败: %v\n", err)
		return
	}
	client := esign.NewClient(testClient.Conf)

	//合同文件列表,通过 go test tests/template_test.go -v -run TestCreateByTemplate 获取
	contractFiles := make([]types.ESignCreateFlowFiles, 0)
	contractFiles = append(contractFiles, types.ESignCreateFlowFiles{
		TemplateId: "213fd9a04b2549818e15871a65bdb41b",
		EFileId:    "14909f52d141499a823c48448560b212",
	})
	contractFiles = append(contractFiles, types.ESignCreateFlowFiles{
		TemplateId: "213fd9a04b2549818e15871a65bdb41b",
		EFileId:    "b3957c0554ba4860aafcb1bd5ea55ac0",
	})
	//设置请求参数
	requestESignCreateFlowData := &types.ESignCreateFlowRequestData{
		SignerName:    signerName,
		SignerPhone:   signerPhone,
		CompanySealID: "", //可以留空,将会使用默认的公司印章
		ContractFiles: contractFiles,
	}
	createFlowResponse, err := client.Sign.ESignCreateFlowOneStep(requestESignCreateFlowData, true)
	if err != nil {
		t.Errorf("Failed to create flow: %v", err)
	}
	//t.Logf("createFlowResponse: %v", utils.JsonMarshalNoEscape(createFlowResponse))

	createFlowResponseData := types.ESignCreateFlowResponseData{}
	err = utils.JsonUnmarshalToStruct(createFlowResponse.Data, &createFlowResponseData)
	if err != nil {
		t.Errorf("ESignCreateFlow JsonUnmarshalToStruct error: %v", err)
		return
	}

	//这里获取到的flowId需要保存起来,可以用来:1.查询签署链接; 2.撤回签署流程; 3.签署完成后查询合同文档;
	t.Logf("createFlowResponseData flowId: %v", createFlowResponseData.ESignFlowId)
}

// TestESignGetExecuteUrlByFlowId 请求e签宝查询签署链接 | go test tests/sign_test.go -v -run TestESignGetExecuteUrlByFlowId
func TestESignGetExecuteUrlByFlowId(t *testing.T) {
	testClient, err := Initialize.NewTestClient()
	if err != nil {
		t.Errorf("创建测试客户端失败: %v\n", err)
		return
	}
	client := esign.NewClient(testClient.Conf)

	//设置请求参数
	flowId := "531daae5b81f47b494d82f0e48603eb6" //TestESignCreateFlowOneStep 中获取到的flowId
	executeUrlResponse, err := client.Sign.GetESignExecuteUrlByFlowId(flowId, signerName, signerPhone, true)
	if err != nil {
		t.Errorf("Failed to get execute url: %v", err)
	}
	//t.Logf("executeUrlResponse: %v", utils.JsonMarshalNoEscape(executeUrlResponse))

	//解析返回结构
	executeUrlResponseData := types.GetExecuteUrlResponseData{}
	err = utils.JsonUnmarshalToStruct(executeUrlResponse.Data, &executeUrlResponseData)
	if err != nil {
		t.Errorf("GetExecuteUrl JsonUnmarshalToStruct error: %v", err)
		return
	}
	executeUrlResponseData.ESignFlowId = flowId

	t.Logf("executeUrlResponseData flowId: %v", executeUrlResponseData.ESignFlowId)
	t.Logf("executeUrlResponseData url: %v", executeUrlResponseData.ESignUrl)
	t.Logf("executeUrlResponseData shortUrl: %v", executeUrlResponseData.ESignShortUrl)
}

// TestESignFlowRevoke 请求e签宝撤回签署流程 | go test tests/sign_test.go -v -run TestESignFlowRevoke
func TestESignFlowRevoke(t *testing.T) {
	testClient, err := Initialize.NewTestClient()
	if err != nil {
		t.Errorf("创建测试客户端失败: %v\n", err)
		return
	}
	client := esign.NewClient(testClient.Conf)

	//设置请求参数
	flowId := "531daae5b81f47b494d82f0e48603eb6" //TestESignCreateFlowOneStep 中获取到的flowId
	revokeResponse, err := client.Sign.ESignFlowRevoke(flowId, true)
	if err != nil {
		t.Errorf("Failed to revoke flow: %v", err)
	}
	t.Logf("revokeResponse: %v", utils.JsonMarshalNoEscape(revokeResponse))
}
