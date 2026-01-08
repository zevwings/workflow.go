package common

// NavigationHandler 导航处理器
// 提供通用的上下箭头键导航逻辑，支持边界处理
type NavigationHandler struct {
	itemCount int
	cyclic    bool // 是否循环导航（暂未实现，预留接口）
}

// NewNavigationHandler 创建导航处理器
//
// 参数:
//   - itemCount: 选项总数
//   - cyclic: 是否循环导航（当到达边界时是否循环到另一端）
//
// 返回:
//   - *NavigationHandler: 导航处理器实例
func NewNavigationHandler(itemCount int, cyclic bool) *NavigationHandler {
	return &NavigationHandler{
		itemCount: itemCount,
		cyclic:    cyclic,
	}
}

// ProcessArrowKey 处理箭头键
// 根据方向（"up" 或 "down"）计算新的索引位置
//
// 参数:
//   - currentIndex: 当前索引位置
//   - direction: 方向（"up" 或 "down"）
//
// 返回:
//   - newIndex: 新的索引位置
//   - shouldRender: 是否需要重新渲染（索引发生变化时返回 true）
func (h *NavigationHandler) ProcessArrowKey(currentIndex int, direction string) (newIndex int, shouldRender bool) {
	if h.itemCount == 0 {
		return currentIndex, false
	}

	if direction == "up" {
		if currentIndex > 0 {
			return currentIndex - 1, true
		}
		// 如果启用循环导航，到达顶部时循环到底部
		if h.cyclic {
			return h.itemCount - 1, true
		}
		return currentIndex, false
	}

	if direction == "down" {
		if currentIndex < h.itemCount-1 {
			return currentIndex + 1, true
		}
		// 如果启用循环导航，到达底部时循环到顶部
		if h.cyclic {
			return 0, true
		}
		return currentIndex, false
	}

	return currentIndex, false
}

// ValidateIndex 验证索引是否有效
// 如果索引无效，返回有效的默认索引（0）
//
// 参数:
//   - index: 要验证的索引
//
// 返回:
//   - validIndex: 有效的索引
func (h *NavigationHandler) ValidateIndex(index int) int {
	if index < 0 || index >= h.itemCount {
		return 0
	}
	return index
}

