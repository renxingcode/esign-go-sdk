package Initialize

import (
	"context"
	"fmt"
	"github.com/renxingcode/esign-go-sdk/config"
	"log"
	"os"
)

type TestClient struct {
	Conf *config.Config
	Ctx  context.Context
}

func NewTestClient() (*TestClient, error) {
	rootPath, err := config.GetProjectRootPath()
	if err != nil {
		return nil, fmt.Errorf("获取项目根目录失败: %w\n", err)
	}
	config.LoadEnvData(rootPath)

	appID := os.Getenv("ESIGN_APP_ID")
	appSecret := os.Getenv("ESIGN_APP_SECRET")
	baseURL := os.Getenv("ESIGN_BASE_URL")
	orgID := os.Getenv("ESIGN_ORG_ID")
	grantType := os.Getenv("ESIGN_GRANT_TYPE")
	isWriteLog := os.Getenv("IS_WRITE_LOG")
	if isWriteLog == "" {
		isWriteLog = "false"
	}
	if appID == "" || appSecret == "" || baseURL == "" || orgID == "" || grantType == "" {
		log.Fatal("ESIGN_APP_ID, ESIGN_APP_SECRET, ESIGN_BASE_URL, ESIGN_ORG_ID and ESIGN_GRANT_TYPE environment variables are required")
	}

	conf, err := config.NewConfig(appID, appSecret, baseURL, orgID, grantType, isWriteLog)
	if err != nil {
		return nil, fmt.Errorf("创建配置失败: %w\n", err)
	}
	ctx := context.Background()
	return &TestClient{
		Conf: conf,
		Ctx:  ctx,
	}, nil
}
