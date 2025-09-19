package api

import (
	"errors"
	"fmt"
	"github.com/renxingcode/esign-go-sdk/types"
	"github.com/renxingcode/esign-go-sdk/utils"
)

// 缓存的常量
const (
	ESignAccessTokenKey        = "esign_access_token" //e签宝的access_token的key
	ESignAccessTokenExpireTime = 7200                 //e签宝的access_token的过期时间,单位秒
)

// 请求e签宝的URL的常量
const (
	GetESignTokenPath            = "/v1/oauth2/access_token"               //获取e签宝的token https://open.esign.cn/doc/opendoc/identity_service/szr5s9
	GetESignTemplateDetailPath   = "/v3/sign-templates/detail"             //获取e签宝的合同模板详情 https://open.esign.cn/doc/opendoc/file-and-template3/pfzut7ho9obc7c5r
	CreateESignFileByTemplate    = "/v1/files/createByTemplate"            //通过模板创建文件 https://open.esign.cn/doc/opendoc/saas_api/cz9d65_sh823i
	GetPersonsIdentityInfo       = "/v3/persons/identity-info"             //查询个人认证信息 https://open.esign.cn/doc/opendoc/auth3/vssvtu
	CreatePersonsIdentity        = "/v1/accounts/createByThirdPartyUserId" //创建个人签署账号 https://open.esign.cn/doc/opendoc/paas_api/ox6nog
	CreateFlowOneStep            = "/api/v2/signflows/createFlowOneStep"   //一步发起签署 https://open.esign.cn/doc/opendoc/paas_api/pwd6l4
	GetESignExecuteUrlByFlowId   = "/v1/signflows/{FLOW_ID}/executeUrl"    //查询签署链接 https://open.esign.cn/doc/opendoc/saas_api/fh3gh1_dwz08n
	ESignFlowRevokePath          = "/v1/signflows/{FLOW_ID}/revoke"        //撤回签署流程 https://open.esign.cn/doc/opendoc/saas_api/hv1dii_uqoamg
	GetESignDocumentsUrlByFlowId = "/v1/signflows/{FLOW_ID}/documents"     //获取签署完成后的文档链接 https://open.esign.cn/doc/opendoc/saas_api/oyqsoq_zknh6g
)

// e签宝返回的code
const (
	ESignResponseCodeSuccess = 0 //成功
)

// GetESignCommonResponse 获取e签宝的通用响应,并解析到结构体
func GetESignCommonResponse(responseJson string) (*types.ESignCommonResponse, error) {
	var responseStruct types.ESignCommonResponse
	err := utils.JsonUnmarshalToStruct(responseJson, &responseStruct)
	if err != nil {
		return nil, err
	}
	return &responseStruct, nil
}

// 构建请求e签宝的headers失败的错误封装
func BuildRequestESignHeadersError(actionName string, err error) error {
	return errors.New(actionName + "构建请求e签宝的headers失败:" + err.Error())
}

// 发送http请求失败的错误封装
func SendHttpRequestError(actionName string, err error) error {
	return errors.New(actionName + "发送http请求失败:" + err.Error())
}

// 解析e签宝响应体失败的错误封装
func ParseESignResponseError(actionName string, err error) error {
	return errors.New(actionName + "解析e签宝响应体失败:" + err.Error())
}

// e签宝返回错误封装
func GetESignResponseError(eSignResponse *types.ESignCommonResponse) error {
	return errors.New(fmt.Sprintf("[e签宝返回错误]code:%d,message:%s", eSignResponse.Code, eSignResponse.Message))
}

// 解析e签宝响应体Data字段失败的错误封装
func ParseESignResponseDataError(actionName string, err error) error {
	return errors.New(actionName + "解析e签宝响应体Data字段失败:" + err.Error())
}
