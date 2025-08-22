package api

import (
	"github.com/renxingcode/esign-go-sdk/types"
	"github.com/renxingcode/esign-go-sdk/utils"
)

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
