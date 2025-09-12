package tests

import (
	"github.com/renxingcode/esign-go-sdk"
	"github.com/renxingcode/esign-go-sdk/Initialize"
	"github.com/renxingcode/esign-go-sdk/types"
	"github.com/renxingcode/esign-go-sdk/utils"
	"testing"
)

// TestESignCreateFlowOneStep 请求e签宝发起签署流程 | go test tests/sign_test.go -v -run TestESignCreateFlowOneStep
func TestESignCreateFlowOneStep(t *testing.T) {
	testClient, err := Initialize.NewTestClient()
	if err != nil {
		t.Errorf("创建测试客户端失败: %v\n", err)
		return
	}
	client := esign.NewClient(testClient.Conf)

	//设置请求参数
	requestESignCreateFlowData := &types.ESignCreateFlowRequestData{
		SignerName:    "张三",
		SignerPhone:   "13945618971",
		CompanySealID: "123456",
		ContractFiles: nil, //todo 合同文件列表
	}
	createFlowResponse, err := client.Sign.ESignCreateFlowOneStep(requestESignCreateFlowData, true)
	if err != nil {
		t.Errorf("Failed to create flow: %v", err)
	}
	t.Logf("createFlowResponse: %v", utils.JsonMarshalNoEscape(createFlowResponse))
}
