package utils

import (
	"encoding/json"
	"fmt"
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
	default:
		return fmt.Errorf("unexpected data type: %T", data)
	}
	return nil
}
