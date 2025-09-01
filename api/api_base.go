package api

import (
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
	GetESignTokenPath = "/v1/oauth2/access_token" //获取e签宝的token
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
