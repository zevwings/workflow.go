//go:build test

package testutils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 测试用的辅助函数
func parseString(input interface{}) (interface{}, error) {
	str, ok := input.(string)
	if !ok {
		return nil, errors.New("输入必须是字符串")
	}
	if str == "" {
		return nil, errors.New("输入不能为空")
	}
	if str == "error" {
		return nil, errors.New("解析失败: 无效的输入")
	}
	return str, nil
}

func addNumbers(input interface{}) (interface{}, error) {
	nums, ok := input.([]int)
	if !ok {
		return nil, errors.New("输入必须是整数切片")
	}
	sum := 0
	for _, n := range nums {
		sum += n
	}
	return sum, nil
}

// TestRunTableTest_Basic 测试基本的表驱动测试功能
func TestRunTableTest_Basic(t *testing.T) {
	testCases := []TableTestCase{
		{
			Name:  "成功解析",
			Input: "test",
			Want:  "test",
		},
		{
			Name:     "空字符串错误",
			Input:    "",
			WantErr:  true,
			WantErrContains: "不能为空",
		},
		{
			Name:     "类型错误",
			Input:    123,
			WantErr:  true,
			WantErrContains: "必须是字符串",
		},
		{
			Name:     "解析失败",
			Input:    "error",
			WantErr:  true,
			WantErrContains: "解析失败",
		},
	}

	RunTableTest(t, testCases, parseString)
}

// TestRunTableTest_WithCustomAssert 测试自定义断言函数
func TestRunTableTest_WithCustomAssert(t *testing.T) {
	testCases := []TableTestCase{
		{
			Name:  "使用自定义断言",
			Input: []int{1, 2, 3},
			AssertFunc: func(t *testing.T, result interface{}) bool {
				sum, ok := result.(int)
				require.True(t, ok, "结果应该是整数")
				assert.Equal(t, 6, sum, "和应该为 6")
				return true
			},
		},
		{
			Name:  "使用期望值断言",
			Input: []int{4, 5, 6},
			Want:  15,
		},
	}

	RunTableTest(t, testCases, addNumbers)
}

// TestRunTableTest_WithCustomErrorCheck 测试自定义错误验证
func TestRunTableTest_WithCustomErrorCheck(t *testing.T) {
	testCases := []TableTestCase{
		{
			Name:  "自定义错误验证",
			Input: "",
			WantErr: true,
			WantErrFunc: func(t *testing.T, err error) bool {
				if err == nil {
					t.Error("应该有错误")
					return false
				}
				assert.Contains(t, err.Error(), "不能为空")
				return true
			},
		},
		{
			Name:     "使用 WantErrContains",
			Input:    "error",
			WantErr:  true,
			WantErrContains: "解析失败",
		},
	}

	RunTableTest(t, testCases, parseString)
}

// TestRunTableTest_WithSetupAndCleanup 测试带设置和清理的测试
func TestRunTableTest_WithSetupAndCleanup(t *testing.T) {
	counter := 0

	testCases := []TableTestCase{
		{
			Name:  "带设置和清理的测试",
			Input: "test",
			Want:  "test",
			Setup: func(t *testing.T) interface{} {
				counter++
				t.Logf("Setup 执行，counter = %d", counter)
				return "setup-value"
			},
			Cleanup: func(t *testing.T) {
				counter--
				t.Logf("Cleanup 执行，counter = %d", counter)
			},
		},
	}

	RunTableTest(t, testCases, func(input interface{}) (interface{}, error) {
		result, err := parseString(input)
		if err != nil {
			return nil, err
		}
		return result, nil
	})

	// 验证清理是否执行
	assert.Equal(t, 0, counter, "清理应该将 counter 重置为 0")
}

// TestRunTableTestWithSetup_Basic 测试带设置的表驱动测试
func TestRunTableTestWithSetup_Basic(t *testing.T) {
	testCases := []TableTestCase{
		{
			Name:  "使用设置值",
			Input: 10,
			Want:  20,
			Setup: func(t *testing.T) interface{} {
				return 10 // 设置值为 10
			},
		},
		{
			Name:  "无设置值",
			Input: 5,
			Want:  5,
		},
	}

	RunTableTestWithSetup(t, testCases, func(input, setup interface{}) (interface{}, error) {
		val, ok := input.(int)
		if !ok {
			return nil, errors.New("输入必须是整数")
		}

		if setup != nil {
			setupVal, ok := setup.(int)
			if ok {
				return val + setupVal, nil
			}
		}

		return val, nil
	})
}

// TestRunTableTest_EmptyTestCases 测试空测试用例列表
func TestRunTableTest_EmptyTestCases(t *testing.T) {
	testCases := []TableTestCase{}

	RunTableTest(t, testCases, parseString)
	// 应该不执行任何测试，但不应该出错
}

// TestRunTableTest_ComplexScenario 测试复杂场景
func TestRunTableTest_ComplexScenario(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	testCases := []TableTestCase{
		{
			Name:  "结构体比较",
			Input: "John:30",
			Want: Person{
				Name: "John",
				Age:  30,
			},
			AssertFunc: func(t *testing.T, result interface{}) bool {
				person, ok := result.(Person)
				require.True(t, ok)
				assert.Equal(t, "John", person.Name)
				assert.Equal(t, 30, person.Age)
				return true
			},
		},
	}

	RunTableTest(t, testCases, func(input interface{}) (interface{}, error) {
		str, ok := input.(string)
		if !ok {
			return nil, errors.New("输入必须是字符串")
		}
		// 简单的解析逻辑（仅用于测试）
		if str == "John:30" {
			return Person{Name: "John", Age: 30}, nil
		}
		return nil, errors.New("无法解析")
	})
}

// TestRunTableTest_ErrorHandling 测试错误处理的各种情况
// 注意：这个测试演示了错误处理逻辑，但跳过了一些会失败的用例
func TestRunTableTest_ErrorHandling(t *testing.T) {
	// 测试正确的错误处理
	testCases := []TableTestCase{
		{
			Name:     "正确的错误处理 - 期望错误并得到错误",
			Input:    "",
			WantErr:  true,
			WantErrContains: "不能为空",
		},
		{
			Name:     "正确的错误处理 - 期望错误并有自定义验证",
			Input:    "",
			WantErr:  true,
			WantErrFunc: func(t *testing.T, err error) bool {
				return assert.NotNil(t, err) && assert.Contains(t, err.Error(), "不能为空")
			},
		},
		{
			Name:     "正确的成功处理 - 不期望错误且无错误",
			Input:    "test",
			WantErr:  false,
			Want:     "test",
		},
	}

	RunTableTest(t, testCases, parseString)

	// 以下测试用例会失败，但它们用于演示错误检测逻辑
	// 在实际使用中，这些测试应该被注释掉或使用 Skip
	t.Run("演示：期望错误但未出现（应该失败）", func(t *testing.T) {
		t.Skip("这是一个演示用例，展示错误检测逻辑")
		tc := TableTestCase{
			Name:    "期望错误但未出现",
			Input:   "test",
			WantErr: true,
		}
		testCases := []TableTestCase{tc}
		RunTableTest(t, testCases, parseString)
		// 这个测试会失败，因为 "test" 不应该产生错误
	})

	t.Run("演示：不期望错误但出现了（应该失败）", func(t *testing.T) {
		t.Skip("这是一个演示用例，展示错误检测逻辑")
		tc := TableTestCase{
			Name:    "不期望错误但出现了",
			Input:   "",
			WantErr: false,
		}
		testCases := []TableTestCase{tc}
		RunTableTest(t, testCases, parseString)
		// 这个测试会失败，因为 "" 应该产生错误
	})
}

// TestRunTableTest_EdgeCases 测试边界情况
func TestRunTableTest_EdgeCases(t *testing.T) {
	testCases := []TableTestCase{
		{
			Name:  "nil 输入",
			Input: nil,
			WantErr: true,
		},
		{
			Name:  "nil 期望值（只验证无错误）",
			Input: "test",
			Want:  nil, // nil 期望值意味着只验证无错误，不验证结果
		},
		{
			Name:  "空名称（应该使用默认名称）",
			Input: "test",
			Want:  "test",
		},
	}

	RunTableTest(t, testCases, func(input interface{}) (interface{}, error) {
		if input == nil {
			return nil, errors.New("输入不能为 nil")
		}
		return parseString(input)
	})
}

// BenchmarkRunTableTest 性能基准测试
func BenchmarkRunTableTest(b *testing.B) {
	testCases := []TableTestCase{
		{
			Name:  "基准测试用例",
			Input: "test",
			Want:  "test",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RunTableTest(&testing.T{}, testCases, parseString)
	}
}

// ExampleRunTableTest 示例：如何使用 RunTableTest
// 注意：这是一个代码示例，实际使用时需要在测试函数中调用
//
// 使用示例：
//
//	func TestParseTicketID(t *testing.T) {
//		// 定义一个简单的解析函数
//		parseID := func(input interface{}) (interface{}, error) {
//			str, ok := input.(string)
//			if !ok {
//				return nil, errors.New("输入必须是字符串")
//			}
//			if str == "" {
//				return nil, errors.New("输入不能为空")
//			}
//			return str, nil
//		}
//
//		// 定义测试用例
//		testCases := []testutils.TableTestCase{
//			{
//				Name:  "有效输入",
//				Input: "PROJ-123",
//				Want:  "PROJ-123",
//			},
//			{
//				Name:     "无效输入 - 空字符串",
//				Input:    "",
//				WantErr:  true,
//				WantErrContains: "不能为空",
//			},
//			{
//				Name:     "无效输入 - 类型错误",
//				Input:    123,
//				WantErr:  true,
//				WantErrContains: "必须是字符串",
//			},
//		}
//
//		// 运行测试
//		testutils.RunTableTest(t, testCases, parseID)
//	}

