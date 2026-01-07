package api

import (
	"context"
	"fmt"

	"github.com/andygrunwald/go-jira/v2/cloud"
)

// UserAPI 提供 User 相关的 REST API 方法
type UserAPI struct {
	client *cloud.Client
	ctx    context.Context
}

// NewUserAPI 创建新的 User API 实例
func NewUserAPI(client *cloud.Client, ctx context.Context) *UserAPI {
	return &UserAPI{
		client: client,
		ctx:    ctx,
	}
}

// GetCurrentUser 获取当前登录用户信息
//
// 返回:
//   - *cloud.User: 当前用户信息
//   - error: 如果获取失败，返回错误
func (api *UserAPI) GetCurrentUser() (*cloud.User, error) {
	user, _, err := api.client.User.GetCurrentUser(api.ctx)
	if err != nil {
		return nil, fmt.Errorf("获取当前用户信息失败: %w", err)
	}

	return user, nil
}

// GetUser 根据 Account ID 获取用户信息
//
// 参数:
//   - accountID: 用户 Account ID
//
// 返回:
//   - *cloud.User: 用户信息
//   - error: 如果获取失败，返回错误
func (api *UserAPI) GetUser(accountID string) (*cloud.User, error) {
	user, _, err := api.client.User.Get(api.ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("获取用户 %s 信息失败: %w", accountID, err)
	}

	return user, nil
}

// FindUsers 搜索用户
//
// 参数:
//   - query: 搜索关键词（用户名或邮箱）
//
// 返回:
//   - []*cloud.User: 用户列表
//   - error: 如果搜索失败，返回错误
func (api *UserAPI) FindUsers(query string) ([]*cloud.User, error) {
	users, _, err := api.client.User.Find(api.ctx, query)
	if err != nil {
		return nil, fmt.Errorf("搜索用户 %s 失败: %w", query, err)
	}

	// 转换 []User 为 []*User
	result := make([]*cloud.User, len(users))
	for i := range users {
		result[i] = &users[i]
	}

	return result, nil
}

