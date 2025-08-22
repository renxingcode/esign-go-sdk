package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"path/filepath"
	"strconv"
)

type Config struct {
	AppID      string
	AppSecret  string
	BaseURL    string
	OrgId      string
	GrantType  string
	IsWriteLog bool
}

// Option 是用于配置客户端的函数类型，采用选项模式，便于未来扩展
type Option func(*Config)

// NewConfig 创建一个默认配置
func NewConfig(appID, appSecret, baseURL, orgID, grantType, isWriteLog string, opts ...Option) (*Config, error) {
	isWriteLogBool, err := strconv.ParseBool(isWriteLog)
	if err != nil {
		fmt.Printf(".env配置文件中的isWriteLog取值异常,未能成功转换为bool类型:%v,应该是true或者false,如果不需要写入日志,可以删除此配置项\n", err)
		return nil, fmt.Errorf("isWriteLog配置异常:%w\n", err)
	}

	conf := &Config{
		AppID:      appID,
		AppSecret:  appSecret,
		BaseURL:    baseURL,
		OrgId:      orgID,
		GrantType:  grantType,
		IsWriteLog: isWriteLogBool,
	}
	for _, opt := range opts {
		opt(conf)
	}
	return conf, nil
}

// LoadEnvData 加载环境变量
func LoadEnvData(path string) {
	// 加载 .env 文件
	if err := godotenv.Load(filepath.Join(path, ".env")); err != nil {
		logx.Errorf("Error loading .env file: %v", err)
	}
}

// GetProjectRootPath 获取项目根目录路径
func GetProjectRootPath() (string, error) {
	// 获取当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// 从当前目录向上查找 go.mod 文件
	for {
		// 检查当前目录是否存在 go.mod 文件
		if _, err := os.Stat(filepath.Join(wd, "go.mod")); err == nil {
			return wd, nil
		}

		// 获取父目录
		parent := filepath.Dir(wd)
		if parent == wd {
			return "", fmt.Errorf("未找到项目根目录")
		}
		wd = parent
	}
}
