package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"io/ioutil"
	"net/http"
)

// SendHttpGetRequest 发送GET请求
func SendHttpGetRequest(requestUrl string, requestHeaders map[string]string, isWriteLog bool) (string, error) {
	// 创建一个新的 GET 请求
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return "", errors.New("创建GET请求失败:" + err.Error())
	}

	// 设置请求头
	if requestHeaders != nil {
		for key, value := range requestHeaders {
			req.Header.Set(key, value)
		}
	}

	// 使用默认的 HTTP 客户端发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logx.Errorf("SendHttpGetRequest发送请求失败: %v", err)
		return "", errors.New("发送GET请求失败:" + err.Error())
	}
	defer resp.Body.Close()

	// 读取响应内容
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	respBodyStr := string(respBody)

	if isWriteLog {
		logData := map[string]interface{}{
			"requestUrl":       requestUrl,
			"requestHeaders":   requestHeaders,
			"responseHttpCode": resp.StatusCode, //HTTP 状态码
			"responseData":     respBodyStr,
		}
		LogxInfow(logData, "SendHttpGetRequestLog")
	}
	return respBodyStr, nil
}

// SendHttpPutRequest 发送PUT请求
func SendHttpPutRequest(requestUrl string, requestHeaders map[string]string, isWriteLog bool) (string, error) {
	// 创建一个新的 GET 请求
	req, err := http.NewRequest("PUT", requestUrl, nil)
	if err != nil {
		return "", errors.New("创建PUT请求失败:" + err.Error())
	}

	// 设置请求头
	if requestHeaders != nil {
		for key, value := range requestHeaders {
			req.Header.Set(key, value)
		}
	}

	// 使用默认的 HTTP 客户端发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logx.Errorf("SendHttpPutRequest发送请求失败: %v", err)
		return "", errors.New("发送PUT请求失败:" + err.Error())
	}
	defer resp.Body.Close()

	// 读取响应内容
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	respBodyStr := string(respBody)

	if isWriteLog {
		logData := map[string]interface{}{
			"requestUrl":       requestUrl,
			"requestHeaders":   requestHeaders,
			"responseHttpCode": resp.StatusCode, //HTTP 状态码
			"responseData":     respBodyStr,
		}
		LogxInfow(logData, "SendHttpPutRequestLog")
	}
	return respBodyStr, nil
}

// SendHttpPostRequest 发送POST请求
func SendHttpPostRequest(requestUrl string, requestBody interface{}, requestHeaders map[string]string, isWriteLog bool) (string, error) {
	// 将请求体转换为 JSON
	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return "", errors.New("请求体转换为JSON失败:" + err.Error())
	}

	// 创建一个新的 POST 请求
	req, err := http.NewRequest("POST", requestUrl, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", errors.New("创建POST请求失败:" + err.Error())
	}

	// 设置请求头
	if requestHeaders != nil {
		for key, value := range requestHeaders {
			req.Header.Set(key, value)
		}
	}

	// 使用默认的 HTTP 客户端发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logx.Errorf("SendHttpPostRequest发送请求失败: %v", err)
		return "", errors.New("发送POST请求失败:" + err.Error())
	}
	defer resp.Body.Close()

	// 读取响应内容
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logx.Errorf("SendHttpPostRequest读取响应内容失败: %v", err)
		return "", errors.New("读取响应内容失败:" + err.Error())
	}
	respBodyStr := string(respBody)

	if isWriteLog {
		logData := map[string]interface{}{
			"requestUrl":       requestUrl,
			"requestHeaders":   requestHeaders,
			"requestBody":      requestBody,
			"responseHttpCode": resp.StatusCode, //HTTP 状态码
			"responseData":     respBodyStr,
		}
		LogxInfow(logData, "SendHttpPostRequestLog")
	}
	return respBodyStr, nil
}
