package tests

import (
	"github.com/renxingcode/esign-go-sdk"
	"github.com/renxingcode/esign-go-sdk/initialize"
	"github.com/renxingcode/esign-go-sdk/types"
	"github.com/renxingcode/esign-go-sdk/utils"
	"testing"
)

// TestESignCreateFlowOneStep 请求e签宝发起签署流程 | go test tests/sign_test.go -v -run TestESignCreateFlowOneStep
func TestESignCreateFlowOneStep(t *testing.T) {
	testClient, err := initialize.NewTestClient()
	if err != nil {
		t.Errorf("创建测试客户端失败: %v\n", err)
		return
	}
	client := esign.NewClient(testClient.Conf)

	//这里模拟一次性签署3份合同,这三份合同假设来自于同一个合同模板
	//合同文件列表,其中的 EFileId 通过 go test tests/template_test.go -v -run TestCreateByTemplate 获取
	contractFiles := make([]types.ESignCreateFlowFiles, 0)
	testESignTemplateId := testClient.Conf.MoreData["eSignTemplateId"].(string)
	contractFiles = append(contractFiles, types.ESignCreateFlowFiles{
		TemplateId: testESignTemplateId,
		EFileId:    "ec7db8732e0e4b1c805844bcbfd2d37b",
	})
	contractFiles = append(contractFiles, types.ESignCreateFlowFiles{
		TemplateId: testESignTemplateId,
		EFileId:    "4aa04dfc0a6b4662af55074bcb15efef",
	})
	contractFiles = append(contractFiles, types.ESignCreateFlowFiles{
		TemplateId: testESignTemplateId,
		EFileId:    "7e9da56c55004c41bf01476f90d0e42b",
	})

	//设置请求参数
	signerName := testClient.Conf.MoreData["signerName"].(string)
	signerPhone := testClient.Conf.MoreData["signerPhone"].(string)
	requestESignCreateFlowData := &types.ESignCreateFlowRequestData{
		SignerName:    signerName,  // 签署人姓名
		SignerPhone:   signerPhone, // 签署人手机号
		CompanySealID: "",          //可以留空,将会使用默认的公司印章
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

	//下一步,直接查询签署链接
	//设置查询签署链接的请求参数
	flowId := createFlowResponseData.ESignFlowId
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

// TestESignGetExecuteUrlByFlowId 请求e签宝查询签署链接 | go test tests/sign_test.go -v -run TestESignGetExecuteUrlByFlowId
func TestESignGetExecuteUrlByFlowId(t *testing.T) {
	testClient, err := initialize.NewTestClient()
	if err != nil {
		t.Errorf("创建测试客户端失败: %v\n", err)
		return
	}
	client := esign.NewClient(testClient.Conf)

	//设置请求参数
	flowId := "01fe697xxxxx" //TestESignCreateFlowOneStep 中获取到的flowId
	signerName := testClient.Conf.MoreData["signerName"].(string)
	signerPhone := testClient.Conf.MoreData["signerPhone"].(string)
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
	testClient, err := initialize.NewTestClient()
	if err != nil {
		t.Errorf("创建测试客户端失败: %v\n", err)
		return
	}
	client := esign.NewClient(testClient.Conf)

	//设置请求参数
	flowId := "01fe697xxxxx" //TestESignCreateFlowOneStep 中获取到的flowId
	revokeResponse, err := client.Sign.ESignFlowRevoke(flowId, true)
	if err != nil {
		t.Errorf("Failed to revoke flow: %v", err)
	}
	t.Logf("revokeResponse: %v", utils.JsonMarshalNoEscape(revokeResponse))
	//{"code":1437111,"message":"非开启状态不允许撤回流程","data":null} todo 和e签宝技术人员沟通
}

// TestESignGetDocumentsUrlByFlowId 查询e签宝签署完成后的文档链接 | go test tests/sign_test.go -v -run TestESignGetDocumentsUrlByFlowId
func TestESignGetDocumentsUrlByFlowId(t *testing.T) {
	testClient, err := initialize.NewTestClient()
	if err != nil {
		t.Errorf("创建测试客户端失败: %v\n", err)
		return
	}
	client := esign.NewClient(testClient.Conf)

	//设置请求参数
	flowId := "01fe697xxxxx" //TestESignCreateFlowOneStep 中获取到的flowId
	eSignDocumentsDocs, err := client.Sign.GetESignDocumentsUrlByFlowId(flowId, true)
	if err != nil {
		t.Errorf("Failed to get documents url: %v", err)
	}

	//todo 对 eSignDocumentsDocs 循环处理,分别上传fileUrl到你自己的OSS服务器
	type GetDocumentsUrlResponseDataDocsForMe struct {
		types.GetDocumentsUrlResponseDataDocs
		OssFileUrl string `json:"ossFileUrl"`
	}

	documentsUrlResponseDataDocsForMe := make([]GetDocumentsUrlResponseDataDocsForMe, 0, len(eSignDocumentsDocs))
	for _, doc := range eSignDocumentsDocs {
		docForMe := GetDocumentsUrlResponseDataDocsForMe{
			GetDocumentsUrlResponseDataDocs: doc,
			OssFileUrl:                      "todo-upload-to-your-oss-server",
		}
		documentsUrlResponseDataDocsForMe = append(documentsUrlResponseDataDocsForMe, docForMe)
	}
	t.Logf("documentsUrlResponseDataDocsForMe: %v", utils.JsonMarshalNoEscape(documentsUrlResponseDataDocsForMe))
}
