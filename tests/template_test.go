package tests

import (
	"github.com/renxingcode/esign-go-sdk"
	"github.com/renxingcode/esign-go-sdk/Initialize"
	"testing"
)

// TestGetTemplateList 测试获取模板列表 | go test tests/template_test.go -v -run TestGetTemplateList
func TestGetTemplateList(t *testing.T) {
	testClient := Initialize.NewTestClient()
	client := esign.NewClient(testClient.Cfg)
	templateList, err := client.Template.GetTemplateList(testClient.Ctx)
	if err != nil {
		t.Errorf("Failed to get template list: %v", err)
	}
	t.Logf("templateList: %v", templateList)
}
