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
		SignerName:    "张三",
		SignerPhone:   "13200010002",
		CompanySealID: "", //可以留空,将会使用默认的公司印章
		ContractFiles: contractFiles,
	}
	createFlowResponse, err := client.Sign.ESignCreateFlowOneStep(requestESignCreateFlowData, true)
	if err != nil {
		t.Errorf("Failed to create flow: %v", err)
	}
	t.Logf("createFlowResponse: %v", utils.JsonMarshalNoEscape(createFlowResponse))
}
