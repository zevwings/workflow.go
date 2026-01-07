package http

// Authorization Basic Authentication 认证信息
//
// 用于 HTTP Basic Authentication 的用户名和密码。
type Authorization struct {
	// Username 用户名（通常是邮箱地址）
	Username string
	// Password 密码（通常是 API token）
	Password string
}

// NewAuthorization 创建新的 Authorization
//
// 创建 Basic Authentication 认证信息。
//
// 参数:
//   - username: 用户名（通常是邮箱地址）
//   - password: 密码（通常是 API token）
//
// 返回:
//   - Authorization 结构体
func NewAuthorization(username, password string) *Authorization {
	return &Authorization{
		Username: username,
		Password: password,
	}
}
