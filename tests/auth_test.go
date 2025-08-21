package tests

import (
	"github.com/renxingcode/esign-go-sdk"
	"github.com/renxingcode/esign-go-sdk/Initialize"
	"testing"
)

// TestGetToken 测试获取Token | go test tests/auth_test.go -v -run TestGetToken
func TestGetToken(t *testing.T) {
	testClient := Initialize.NewTestClient()
	client := esign.NewClient(testClient.Cfg)
	token, err := client.Auth.GetToken(testClient.Ctx)
	if err != nil {
		t.Errorf("Failed to get token: %v", err)
	}
	t.Logf("token: %v", token)
}
