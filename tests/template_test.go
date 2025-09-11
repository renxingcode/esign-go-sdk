package tests

import (
	"github.com/renxingcode/esign-go-sdk"
	"github.com/renxingcode/esign-go-sdk/Initialize"
	"github.com/renxingcode/esign-go-sdk/api"
	"github.com/renxingcode/esign-go-sdk/types"
	"github.com/renxingcode/esign-go-sdk/utils"
	"testing"
)

// TestGetESignTemplateDetail 测试获取流程模版详情 | go test tests/template_test.go -v -run TestGetESignTemplateDetail
func TestGetESignTemplateDetail(t *testing.T) {
	testClient, err := Initialize.NewTestClient()
	if err != nil {
		t.Errorf("创建测试客户端失败: %v\n", err)
		return
	}
	client := esign.NewClient(testClient.Conf)

	eSignTemplateId := "213fd9a04b2549818e15871a65bdb41b" //流程模板ID
	queryComponents := true                               //是否需要查询控件信息:true-是 false-否(默认false)
	writeLog := false                                     //获取模板返回的数据量很大,因此可以根据情况考虑是否关闭写入日志
	templateDetail, err := client.Template.GetESignTemplateDetail(eSignTemplateId, queryComponents, writeLog)
	if err != nil {
		t.Errorf("Failed to get template detail: %v", err)
	}
	t.Logf("templateDetail: %v", utils.JsonMarshalNoEscape(templateDetail))
}

// TestCreateByTemplate 测试通过模板创建文件 | go test tests/template_test.go -v -run TestCreateByTemplate
func TestCreateByTemplate(t *testing.T) {
	testClient, err := Initialize.NewTestClient()
	if err != nil {
		t.Errorf("创建测试客户端失败: %v\n", err)
		return
	}
	client := esign.NewClient(testClient.Conf)

	// 1. 获取模板信息
	eSignTemplateId := "213fd9a04b2549818e15871a65bdb41b" //流程模板ID
	queryComponents := true                               //是否需要查询控件信息:true-是 false-否(默认false)
	writeLog := true                                      //获取模板返回的数据量很大,因此可以根据情况考虑是否关闭写入日志
	eSignResponseTemplateDetail, err := client.Template.GetESignTemplateDetail(eSignTemplateId, queryComponents, writeLog)
	if err != nil {
		t.Errorf("GetESignTemplateDetail error1: %v", err)
		return
	}
	//utils.LogxInfow(eSignResponseTemplateDetail, "eSignResponseTemplateDetail")

	// 2. 判断e签宝返回的code是否成功
	if eSignResponseTemplateDetail.Code != api.ESignResponseCodeSuccess {
		t.Errorf("GetESignTemplateDetail error2: %v", api.GetESignResponseError(eSignResponseTemplateDetail))
		return
	}

	// 3. 解析模板数据到结构体
	eSignTemplateData := types.GetESignTemplateDetailResponse{}
	err = utils.JsonUnmarshalToStruct(eSignResponseTemplateDetail.Data, &eSignTemplateData)
	if err != nil {
		t.Errorf("GetESignTemplateDetail JsonUnmarshalToStruct error: %v", err)
		return
	}

	// 4. 生成合同模板数据
	eSignTemplateDocFileId := eSignTemplateData.Docs[0].FileId
	eSignTemplateDocFileName := "rxESignGoSdkDemo-" + utils.GetCurrentTime() + "-" + eSignTemplateData.Docs[0].FileName //自定义生成后的模板文件名称 todo 改为你自己的格式
	eSignTemplateDocSimpleFormFields := make(map[string]string)
	for _, field := range eSignTemplateData.Participants {
		for _, component := range field.Components {
			//这里先简单粗暴的把控件的名称赋值给控件编码,正常业务场景下应该根据你自己的业务去赋值 todo
			eSignTemplateDocSimpleFormFields[component.ComponentKey] = component.ComponentName
		}
	}

	// 5. 请求e签宝创建合同模板数据CreateESignFileByTemplateRequest
	createByTemplateResponse, err := client.Template.CreateByTemplate(eSignTemplateDocFileId, eSignTemplateDocFileName, eSignTemplateDocSimpleFormFields, true)
	if err != nil {
		t.Errorf("CreateByTemplate error1: %v", err)
		return
	}
	utils.LogxInfow(createByTemplateResponse, "createByTemplateResponse")
}
