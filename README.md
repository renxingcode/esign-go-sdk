# Go语言的e签宝SDK

使用Go语言对接e签宝签署合同的SDK

## 项目简介

这是一个基于Go语言开发的e签宝SDK，用于对接e签宝电子合同签署平台API，实现电子合同的创建、签署、查询等功能。该SDK采用模块化设计，结构清晰，易于集成和扩展。

## 目录结构

```text
├── api/               # API服务实现
│   ├── account_api/   # 账户相关API
│   ├── auth_api/      # 认证相关API
│   ├── sign_api/      # 签署流程API
│   ├── template_api/  # 模板相关API
│   └── api_base.go    # API基础常量和工具
├── config/            # 配置相关代码
│   ├── esign_config.go  # e签宝配置
│   └── redis_config.go  # Redis缓存配置
├── Initialize/        # 初始化工具
│   └── tests_init.go   # 测试环境初始化
├── types/             # 数据结构定义
│   └── types.go        # 所有API请求/响应结构体
├── utils/             # 通用工具函数
│   └── utils.go        # JSON处理、日志等工具函数
├── tests/             # 测试用例
├── .env               # 环境变量配置文件
├── .env.demo          # 环境变量配置示例
├── client.go          # SDK主入口
├── go.mod             # Go模块依赖
└── go.sum             # 依赖版本锁定
```

## 安装方法

```bash
# 使用go mod安装
go get github.com/renxingcode/esign-go-sdk
```

## 快速开始

### 1. 配置环境变量

复制.env.demo文件为.env，并填写您的e签宝账号信息：

```bash
cp .env.demo .env
# 编辑.env文件，填写您的配置
```

.env文件配置项说明：

```text
ESIGN_APP_ID=您的AppID
ESIGN_APP_SECRET=您的AppSecret
ESIGN_BASE_URL=e签宝API地址（如：https://smlopenapi.esign.cn）
ESIGN_ORG_ID=您的机构账号ID
ESIGN_GRANT_TYPE=授权类型（通常为client_credentials）
IS_WRITE_LOG=true/false（是否记录日志）
```

### 2. 初始化SDK客户端

```go
import (
"github.com/renxingcode/esign-go-sdk"
"github.com/renxingcode/esign-go-sdk/config"
)

// 方法一：通过环境变量初始化
conf, err := config.LoadConfigFromEnv()
if err != nil {
// 处理错误
}
client := esign.NewClient(conf)

// 方法二：直接传入配置参数
conf, err := config.NewConfig("您的AppID", "您的AppSecret", "e签宝API地址", "您的机构账号ID", "授权类型", "是否记录日志")
if err != nil {
// 处理错误
}
client := esign.NewClient(conf)
```

## 核心功能模块

### 1. 认证管理

SDK会自动管理e签宝的访问令牌，包括获取、缓存和刷新：

```go
// 获取访问令牌（通常不需要手动调用，SDK内部会自动处理）
token, err := client.Auth.GetESignToken(true) // true表示使用缓存
if err != nil {
// 处理错误
}
```

### 2. 模板管理

```go
// 获取模板详情
templateDetail, err := client.Template.GetESignTemplateDetail("模板ID", true, true)
if err != nil {
// 处理错误
}

// 通过模板创建文件
simpleFormFields := map[string]string{
"field1": "value1",
"field2": "value2",
}
fileResponse, err := client.Template.CreateByTemplate("模板ID", "文件名称", simpleFormFields, true)
if err != nil {
// 处理错误
}
```

### 3. 账户管理

```go
// 查询个人认证信息
accountInfo, err := client.Account.GetESignPersonsIdentityInfo("手机号或邮箱", true)
if err != nil {
// 处理错误
}
```

### 4. 签署流程管理

```go
import "github.com/renxingcode/esign-go-sdk/types"

// 一步发起签署流程
contractFiles := make([]types.ESignCreateFlowFiles, 0)
contractFiles = append(contractFiles, types.ESignCreateFlowFiles{
TemplateId: "模板ID",
EFileId:    "文件ID",
})

requestData := &types.ESignCreateFlowRequestData{
SignerName:    "签署人姓名",
SignerPhone:   "签署人手机号",
CompanySealID: "公司印章ID（可选）",
ContractFiles: contractFiles,
}

flowResponse, err := client.Sign.ESignCreateFlowOneStep(requestData, true)
if err != nil {
// 处理错误
}

// 查询签署链接
flowId := "签署流程ID"
executeUrlResponse, err := client.Sign.GetESignExecuteUrlByFlowId(flowId, "签署人姓名", "签署人手机号", true)
if err != nil {
// 处理错误
}

// 撤回签署流程
revokeResponse, err := client.Sign.ESignFlowRevoke(flowId, true)
if err != nil {
// 处理错误
}
```

## 运行测试

SDK包含完整的测试用例，可以通过以下命令运行：

```bash
# 运行所有测试
cd /path/to/esign-go-sdk
go test ./tests/... -v

# 运行特定测试
# 获取模板详情测试
go test tests/template_test.go -v -run TestGetESignTemplateDetail

# 发起签署流程测试
go test tests/sign_test.go -v -run TestESignCreateFlowOneStep

# 查询个人认证信息测试
go test tests/account_test.go -v -run TestGetESignPersonsIdentityInfo
```

## 依赖说明

- `github.com/joho/godotenv` - 环境变量加载
- `github.com/redis/go-redis/v9` - Redis客户端，用于缓存token
- `github.com/zeromicro/go-zero` - 提供日志等工具

## 注意事项

- 使用前请确保您已经获取了e签宝的开发者账号和相关配置信息
- 部分功能可能需要开通特定权限，请参考e签宝官方文档
- SDK使用Redis缓存token，默认配置为本地Redis服务
- 在生产环境中，请确保妥善保管您的AppSecret和其他敏感信息

## License

MIT

## 联系方式

如有问题或建议，请联系项目维护者：renxingcode