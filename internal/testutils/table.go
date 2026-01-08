//go:build test

package testutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TableTestCase 表驱动测试用例
// 用于简化表驱动测试的编写
type TableTestCase struct {
	// Name 测试用例名称（必需）
	Name string

	// Input 输入值
	Input interface{}

	// Want 期望的输出值（当 WantErr 为 false 时使用）
	Want interface{}

	// WantErr 是否期望出现错误
	WantErr bool

	// WantErrContains 期望错误信息中包含的字符串
	WantErrContains string

	// WantErrFunc 自定义错误验证函数（优先级高于 WantErrContains）
	// 返回 true 表示验证通过，false 表示验证失败
	WantErrFunc func(*testing.T, error) bool

	// AssertFunc 自定义断言函数（优先级高于 Want）
	// 用于复杂的结果验证，如部分字段匹配、深度比较等
	// 返回 true 表示验证通过，false 表示验证失败
	AssertFunc func(*testing.T, interface{}) bool

	// Setup 测试用例前置设置函数（可选）
	// 返回的值会传递给测试函数和断言函数
	Setup func(*testing.T) interface{}

	// Cleanup 测试用例后置清理函数（可选）
	Cleanup func(*testing.T)
}

// RunTableTest 运行表驱动测试
// testCases: 测试用例列表
// fn: 执行测试的函数，接收输入值，返回结果和错误
//
// 示例：
//
//	testCases := []testutils.TableTestCase{
//	    {
//	        Name:  "有效输入",
//	        Input: "PROJ-123",
//	        Want:  "PROJ-123",
//	    },
//	    {
//	        Name:     "无效输入",
//	        Input:    "invalid",
//	        WantErr:  true,
//	        WantErrContains: "invalid format",
//	    },
//	}
//
//	testutils.RunTableTest(t, testCases, func(input interface{}) (interface{}, error) {
//	    return ParseTicketID(input.(string))
//	})
func RunTableTest(t *testing.T, testCases []TableTestCase, fn func(interface{}) (interface{}, error)) {
	t.Helper()

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Helper()

			// 执行前置设置（如果有）
			var setupResult interface{}
			if tc.Setup != nil {
				setupResult = tc.Setup(t)
			}

			// 执行清理（如果有，使用 defer 确保执行）
			if tc.Cleanup != nil {
				defer tc.Cleanup(t)
			}

			// 执行测试函数
			result, err := fn(tc.Input)

			// 错误验证
			if tc.WantErr {
				require.Error(t, err, "期望出现错误，但没有错误")

				// 如果提供了自定义错误验证函数，使用它
				if tc.WantErrFunc != nil {
					if !tc.WantErrFunc(t, err) {
						t.Errorf("自定义错误验证失败")
					}
				} else if tc.WantErrContains != "" {
					// 否则检查错误信息是否包含指定字符串
					assert.Contains(t, err.Error(), tc.WantErrContains, "错误信息应该包含指定字符串")
				}
			} else {
				require.NoError(t, err, "不期望出现错误，但出现了错误: %v", err)

				// 结果验证
				if tc.AssertFunc != nil {
					// 使用自定义断言函数
					if !tc.AssertFunc(t, result) {
						t.Errorf("自定义断言失败")
					}
				} else if tc.Want != nil || setupResult != nil {
					// 如果有期望值或设置结果，进行比较
					expected := tc.Want
					if setupResult != nil && expected == nil {
						// 如果只有设置结果，使用设置结果作为期望值
						expected = setupResult
					}
					assert.Equal(t, expected, result, "结果应该与期望值匹配")
				}
			}
		})
	}
}

// RunTableTestWithSetup 运行带设置的表驱动测试
// 与 RunTableTest 类似，但每个测试用例都有独立的设置
//
// 示例：
//
//	testCases := []testutils.TableTestCase{
//	    {
//	        Name: "测试用例1",
//	        Input: "input1",
//	        Want: "output1",
//	        Setup: func(t *testing.T) interface{} {
//	            // 设置测试环境
//	            return "setup-value"
//	        },
//	        Cleanup: func(t *testing.T) {
//	            // 清理测试环境
//	        },
//	    },
//	}
//
//	testutils.RunTableTestWithSetup(t, testCases, func(input, setup interface{}) (interface{}, error) {
//	    // 可以使用 setup 的值
//	    return process(input.(string), setup.(string))
//	})
func RunTableTestWithSetup(t *testing.T, testCases []TableTestCase, fn func(interface{}, interface{}) (interface{}, error)) {
	t.Helper()

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Helper()

			// 执行前置设置（如果有）
			var setupResult interface{}
			if tc.Setup != nil {
				setupResult = tc.Setup(t)
			}

			// 执行清理（如果有，使用 defer 确保执行）
			if tc.Cleanup != nil {
				defer tc.Cleanup(t)
			}

			// 执行测试函数（传入输入和设置结果）
			result, err := fn(tc.Input, setupResult)

			// 错误验证
			if tc.WantErr {
				require.Error(t, err, "期望出现错误，但没有错误")

				// 如果提供了自定义错误验证函数，使用它
				if tc.WantErrFunc != nil {
					if !tc.WantErrFunc(t, err) {
						t.Errorf("自定义错误验证失败")
					}
				} else if tc.WantErrContains != "" {
					// 否则检查错误信息是否包含指定字符串
					assert.Contains(t, err.Error(), tc.WantErrContains, "错误信息应该包含指定字符串")
				}
			} else {
				require.NoError(t, err, "不期望出现错误，但出现了错误: %v", err)

				// 结果验证
				if tc.AssertFunc != nil {
					// 使用自定义断言函数
					if !tc.AssertFunc(t, result) {
						t.Errorf("自定义断言失败")
					}
				} else if tc.Want != nil {
					assert.Equal(t, tc.Want, result, "结果应该与期望值匹配")
				}
			}
		})
	}
}

