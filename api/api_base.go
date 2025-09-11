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
	GetESignTokenPath          = "/v1/oauth2/access_token"    //获取e签宝的token
	GetESignTemplateDetailPath = "/v3/sign-templates/detail"  //获取e签宝的合同模板详情
	CreateESignFileByTemplate  = "/v1/files/createByTemplate" //通过模板创建文件
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

func GetESignResponseError(eSignResponse *types.ESignCommonResponse) error {
	return errors.New(fmt.Sprintf("[e签宝返回错误]code:%d,message:%s", eSignResponse.Code, eSignResponse.Message))
}
