package http

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
)

// HttpResponse HTTP 响应格式
//
// 封装 HTTP 响应的状态码、状态文本、响应数据和 Headers。
// 响应体延迟解析，通过方法（AsJSON, AsText 等）来解析。
type HttpResponse struct {
	// Status HTTP 状态码（如 200、404、500）
	Status int
	// StatusText HTTP 状态文本（如 "OK"、"Not Found"、"Internal Server Error"）
	StatusText string
	// Headers HTTP 响应 Headers
	Headers map[string]string
	// bodyBytes 缓存的响应体字节（用于延迟解析）
	bodyBytes []byte
}

// FromRestyResponse 从 resty.Response 创建 HttpResponse
//
// 只提取元数据（status、statusText、headers），并缓存响应体字节。
// 响应体通过后续的方法（AsJSON, AsText 等）来解析。
//
// 参数:
//   - response: resty 的响应对象
//
// 返回:
//   - HttpResponse: 封装后的响应
//   - error: 如果读取响应体失败，返回错误
func FromRestyResponse(response *resty.Response) (*HttpResponse, error) {
	status := response.StatusCode()
	statusText := response.Status()
	headers := make(map[string]string)
	for k, v := range response.Header() {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}

	// 缓存响应体字节（可以多次解析）
	bodyBytes := response.Body()

	return &HttpResponse{
		Status:     status,
		StatusText: statusText,
		Headers:    headers,
		bodyBytes:  bodyBytes,
	}, nil
}

// IsSuccess 检查是否为成功响应（状态码 200-299）
//
// 判断 HTTP 状态码是否在成功范围内（200-299）。
//
// 返回:
//   - bool: 如果状态码在 200-299 范围内返回 true，否则返回 false
func (r *HttpResponse) IsSuccess() bool {
	return r.Status >= 200 && r.Status < 300
}

// IsError 检查是否为错误响应
//
// 判断 HTTP 状态码是否不在成功范围内（即状态码 < 200 或 >= 300）。
//
// 返回:
//   - bool: 如果状态码不在 200-299 范围内返回 true，否则返回 false
func (r *HttpResponse) IsError() bool {
	return !r.IsSuccess()
}

// AsJSON 解析为 JSON（便捷方法）
//
// 将响应体解析为 JSON 并反序列化为类型 T。
//
// 类型参数:
//   - T: 目标类型，必须能够被 JSON 反序列化
//
// 返回:
//   - T: 解析后的数据
//   - error: 如果 JSON 解析失败，返回错误
func AsJSON[T any](r *HttpResponse) (T, error) {
	var result T

	// 直接解析到目标类型，避免双重序列化/反序列化
	if len(r.bodyBytes) == 0 || isWhitespace(r.bodyBytes) {
		return result, nil
	}

	if err := json.Unmarshal(r.bodyBytes, &result); err != nil {
		responseText := string(r.bodyBytes)
		preview := responseText
		if len(responseText) > 200 {
			preview = responseText[:200] + "..."
		}
		return result, fmt.Errorf("failed to parse JSON response (HTTP %d). Response preview: %s: %w", r.Status, preview, err)
	}

	return result, nil
}

// AsText 解析为文本（便捷方法）
//
// 将响应体解析为 UTF-8 文本字符串。
//
// 返回:
//   - string: 响应体的文本内容
//   - error: 如果读取响应体失败或不是有效的 UTF-8，返回错误
func (r *HttpResponse) AsText() (string, error) {
	return ParseText(r.bodyBytes, r.Status)
}

// AsBytes 解析为字节
//
// 返回响应体的原始字节。
//
// 返回:
//   - []byte: 响应体字节
func (r *HttpResponse) AsBytes() []byte {
	return r.bodyBytes
}

// EnsureSuccess 确保响应是成功的，否则返回错误
//
// 检查 HTTP 状态码是否在成功范围内（200-299）。
// 如果响应失败，返回包含状态码和响应体的错误信息。
//
// 返回:
//   - *HttpResponse: 如果响应成功，返回自身
//   - error: 如果响应失败，返回包含错误信息的错误
func (r *HttpResponse) EnsureSuccess() (*HttpResponse, error) {
	if !r.IsSuccess() {
		text, _ := r.AsText()
		if text == "" {
			text = "Unable to read response body"
		}
		return nil, fmt.Errorf("HTTP request failed with status %d: %s", r.Status, text)
	}
	return r, nil
}

// EnsureSuccessWith 确保响应是成功的，使用自定义错误处理器
//
// 检查 HTTP 状态码是否在成功范围内（200-299）。
// 如果响应失败，使用提供的错误处理器生成错误。
// 如果响应成功，返回自身以便链式调用。
//
// 参数:
//   - errorHandler: 错误处理函数，接收 *HttpResponse 并返回错误
//
// 返回:
//   - *HttpResponse: 如果响应成功，返回自身
//   - error: 如果响应失败，返回错误处理器生成的错误
func (r *HttpResponse) EnsureSuccessWith(errorHandler func(*HttpResponse) error) (*HttpResponse, error) {
	if !r.IsSuccess() {
		return nil, errorHandler(r)
	}
	return r, nil
}

// ExtractErrorMessage 提取错误消息（通用方法）
//
// 尝试从响应体中提取错误信息，支持多种常见的错误格式：
// - JSON 格式：尝试从 error.message、error 或 message 字段提取
// - 文本格式：如果无法解析为 JSON，则作为文本返回
//
// 返回:
//   - string: 提取的错误消息字符串。如果无法提取，返回格式化的 JSON 或文本内容。
func (r *HttpResponse) ExtractErrorMessage() string {
	// 尝试解析错误响应为 JSON，提取详细的错误信息
	var errorJSON map[string]interface{}
	if err := json.Unmarshal(r.bodyBytes, &errorJSON); err == nil {
		// 尝试提取常见的错误字段
		var errorDetail string
		if errorObj, ok := errorJSON["error"].(map[string]interface{}); ok {
			if msg, ok := errorObj["message"].(string); ok {
				errorDetail = msg
			}
		}
		if errorDetail == "" {
			if errStr, ok := errorJSON["error"].(string); ok {
				errorDetail = errStr
			}
		}
		if errorDetail == "" {
			if msg, ok := errorJSON["message"].(string); ok {
				errorDetail = msg
			}
		}

		jsonStr, _ := json.Marshal(errorJSON)
		if errorDetail != "" {
			return fmt.Sprintf("%s (details: %s)", string(jsonStr), errorDetail)
		}
		return string(jsonStr)
	}

	// 如果不是 JSON，尝试作为文本解析
	text, err := r.AsText()
	if err != nil {
		// 如果无法解析为文本，返回原始字节的字符串表示
		return string(r.bodyBytes)
	}
	return text
}

// ParseWith 使用指定的 Parser 解析响应（通用方法）
//
// 允许使用自定义的 Parser 来解析响应体。
//
// 参数:
//   - parser: 响应解析器实例
//
// 返回:
//   - interface{}: 解析后的数据
//   - error: 如果解析失败，返回错误
func (r *HttpResponse) ParseWith(parser ResponseParser) (interface{}, error) {
	return parser.Parse(r.bodyBytes, r.Status)
}

// GetHeader 获取指定的 Header 值
//
// 参数:
//   - key: Header 键
//
// 返回:
//   - string: Header 值，如果不存在返回空字符串
//   - bool: 是否存在该 Header
func (r *HttpResponse) GetHeader(key string) (string, bool) {
	// 不区分大小写查找
	keyLower := strings.ToLower(key)
	for k, v := range r.Headers {
		if strings.ToLower(k) == keyLower {
			return v, true
		}
	}
	return "", false
}

