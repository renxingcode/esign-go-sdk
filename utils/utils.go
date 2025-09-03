package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"log"
)

// JsonMarshal 生成紧凑的 JSON 格式，没有换行和缩进
func JsonMarshal(data any) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("编码 JSON 失败: %s", err)
	}
	return string(jsonData)
}

// JsonMarshalIndent 生成格式化的 JSON 格式，包含换行和缩进，便于阅读
func JsonMarshalIndent(data any) string {
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Fatalf("编码 JSON 失败: %s", err)
	}
	return string(jsonData)
}

// JsonMarshalNoEscape 生成紧凑的 JSON 格式，不转义 HTML 字符
func JsonMarshalNoEscape(data any) string {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(data)
	if err != nil {
		log.Fatalf("编码 JSON 失败: %s", err)
	}

	result := buffer.String()
	if len(result) > 0 && result[len(result)-1] == '\n' {
		result = result[:len(result)-1]
	}
	return result
}

// JsonMarshalIndentNoEscape 生成格式化的 JSON 格式，不转义 HTML 字符
func JsonMarshalIndentNoEscape(data any) string {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "    ")
	err := encoder.Encode(data)
	if err != nil {
		log.Fatalf("编码 JSON 失败: %s", err)
	}

	result := buffer.String()
	if len(result) > 0 && result[len(result)-1] == '\n' {
		result = result[:len(result)-1]
	}
	return result
}

// JsonUnmarshalToStruct 将JSON字符串解析为指定的结构体
func JsonUnmarshalToStruct(jsonData any, obj any) error {
	// 先检查jsonData的实际类型
	switch data := jsonData.(type) {
	case string:
		// 如果是字符串，直接 Unmarshal
		if jsonData == "" || jsonData == "[]" || jsonData == "{}" {
			return nil
		}
		err := json.Unmarshal([]byte(data), obj)
		if err != nil {
			return err
		}
	case map[string]any:
		// 如果是map，需要先转换为 JSON 再 Unmarshal
		jsonByte, err := json.Marshal(data)
		if err != nil {
			return err
		}
		err = json.Unmarshal(jsonByte, obj)
	case nil:
		return nil
	default:
		return fmt.Errorf("JsonUnmarshalToStruct: unexpected data type: %T", data)
	}
	return nil
}

// 写入日志方法封装
func LogxInfow(logData any, logTitle string) {
	logx.Infow(JsonMarshalNoEscape(logData),
		logx.Field("logxTitle", logTitle),
		logx.Field("logxDataType", fmt.Sprintf("%T", logData)),
	)
}
func LogxErrorw(logData any, logTitle string) {
	logx.Errorw(JsonMarshalNoEscape(logData),
		logx.Field("logxTitle", logTitle),
		logx.Field("logxDataType", fmt.Sprintf("%T", logData)),
	)
}
