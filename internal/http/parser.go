package http

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ResponseParser 响应解析器接口
//
// 定义如何将响应体字节解析为特定类型。
// 不同的格式（JSON、Text、XML 等）可以实现此接口来提供解析逻辑。
type ResponseParser interface {
	// Parse 解析响应体
	//
	// 参数:
	//   - bytes: 响应体字节
	//   - status: HTTP 状态码（用于错误处理和验证）
	//
	// 返回:
	//   - interface{}: 解析后的数据
	//   - error: 如果解析失败，返回错误
	Parse(bytes []byte, status int) (interface{}, error)
}

// JsonParser JSON 解析器
//
// 将响应体解析为 JSON 格式。
type JsonParser struct{}

// Parse 解析 JSON 响应体
func (p *JsonParser) Parse(bytes []byte, status int) (interface{}, error) {
	// 处理空响应
	if len(bytes) == 0 || isWhitespace(bytes) {
		var result interface{}
		if err := json.Unmarshal([]byte("null"), &result); err != nil {
			return nil, fmt.Errorf("failed to parse empty response as JSON: %w", err)
		}
		return result, nil
	}

	var result interface{}
	if err := json.Unmarshal(bytes, &result); err != nil {
		responseText := string(bytes)
		preview := responseText
		if len(responseText) > 200 {
			preview = responseText[:200] + "..."
		}
		return nil, fmt.Errorf("failed to parse JSON response (HTTP %d). Response preview: %s: %w", status, preview, err)
	}

	return result, nil
}

// ParseJSON 解析 JSON 响应体到指定类型
//
// 参数:
//   - bytes: 响应体字节
//   - status: HTTP 状态码
//   - v: 目标类型指针
//
// 返回:
//   - error: 如果解析失败，返回错误
func ParseJSON(bytes []byte, status int, v interface{}) error {
	parser := &JsonParser{}
	result, err := parser.Parse(bytes, status)
	if err != nil {
		return err
	}

	// 将结果转换为 JSON 再解析到目标类型
	resultBytes, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal parsed result: %w", err)
	}

	return json.Unmarshal(resultBytes, v)
}

// TextParser 文本解析器
//
// 将响应体解析为 UTF-8 文本字符串。
type TextParser struct{}

// Parse 解析文本响应体
func (p *TextParser) Parse(bytes []byte, status int) (interface{}, error) {
	// 检查状态码
	if status < 200 || status >= 300 {
		return nil, fmt.Errorf("HTTP request failed with status %d", status)
	}

	text := string(bytes)
	return text, nil
}

// ParseText 解析文本响应体
//
// 参数:
//   - bytes: 响应体字节
//   - status: HTTP 状态码
//
// 返回:
//   - string: 解析后的文本
//   - error: 如果解析失败，返回错误
func ParseText(bytes []byte, status int) (string, error) {
	parser := &TextParser{}
	result, err := parser.Parse(bytes, status)
	if err != nil {
		return "", err
	}

	text, ok := result.(string)
	if !ok {
		return "", fmt.Errorf("unexpected result type: %T", result)
	}

	return text, nil
}

// isWhitespace 检查字节数组是否只包含空白字符
func isWhitespace(bytes []byte) bool {
	return len(bytes) == 0 || len(strings.TrimSpace(string(bytes))) == 0
}
