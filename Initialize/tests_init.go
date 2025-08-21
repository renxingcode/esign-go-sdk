package Initialize

import (
	"context"
	"github.com/renxingcode/esign-go-sdk/config"
	"log"
	"os"
)

type TestClient struct {
	Cfg *config.Config
	Ctx context.Context
}

func NewTestClient() *TestClient {
	rootPath, err := config.GetProjectRootPath()
	if err != nil {
		log.Fatalf("获取项目根目录失败: %s", err)
	}
	config.LoadEnvData(rootPath)

	appID := os.Getenv("ESIGN_APP_ID")
	appSecret := os.Getenv("ESIGN_APP_SECRET")
	baseURL := os.Getenv("ESIGN_BASE_URL")
	orgID := os.Getenv("ESIGN_ORG_ID")
	if appID == "" || appSecret == "" || baseURL == "" || orgID == "" {
		log.Fatal("ESIGN_APP_ID, ESIGN_APP_SECRET, ESIGN_BASE_URL and ESIGN_ORG_ID environment variables are required")
	}

	cfg := config.NewConfig(appID, appSecret, baseURL, orgID)
	ctx := context.Background()
	return &TestClient{
		Cfg: cfg,
		Ctx: ctx,
	}
}
