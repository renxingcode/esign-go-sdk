package tests

import (
	"github.com/renxingcode/esign-go-sdk"
	"github.com/renxingcode/esign-go-sdk/Initialize"
	"github.com/renxingcode/esign-go-sdk/utils"
	"testing"
)

// TestGetTemplateList 测试获取模板列表 | go test tests/template_test.go -v -run TestGetTemplateList
func TestGetTemplateList(t *testing.T) {
	testClient, err := Initialize.NewTestClient()
	if err != nil {
		t.Errorf("创建测试客户端失败: %v\n", err)
		return
	}
	client := esign.NewClient(testClient.Conf)

	eSignTemplateId := "your_template_id" //流程模板ID
	queryComponents := true               //是否需要查询控件信息:true-是 false-否(默认false)
	writeLog := true                      //获取模板返回的数据量很大,因此可以根据情况考虑是否关闭写入日志
	templateList, err := client.Template.GetESignTemplateList(eSignTemplateId, queryComponents, writeLog)
	if err != nil {
		t.Errorf("Failed to get template list: %v", err)
	}
	t.Logf("templateList: %v", utils.JsonMarshalNoEscape(templateList))
}
