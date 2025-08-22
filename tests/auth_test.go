package tests

import (
	"github.com/renxingcode/esign-go-sdk"
	"github.com/renxingcode/esign-go-sdk/Initialize"
	"testing"
)

// TestGetToken 测试获取Token | go test tests/auth_test.go -v -run TestGetToken
func TestGetToken(t *testing.T) {
	testClient, err := Initialize.NewTestClient()
	if err != nil {
		t.Errorf("创建测试客户端失败: %v\n", err)
		return
	}
	client := esign.NewClient(testClient.Conf)
	token, err := client.Auth.GetESignToken(testClient.Ctx)
	if err != nil {
		t.Errorf("Failed to get token: %v\n", err)
		return
	}
	t.Logf("token: %v", token)
}
