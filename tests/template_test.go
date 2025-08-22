package tests

import (
	"github.com/renxingcode/esign-go-sdk"
	"github.com/renxingcode/esign-go-sdk/Initialize"
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
	templateList, err := client.Template.GetESignTemplateList(testClient.Ctx)
	if err != nil {
		t.Errorf("Failed to get template list: %v", err)
	}
	t.Logf("templateList: %v", templateList)
}
