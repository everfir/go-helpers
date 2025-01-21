package slice

import (
	"cmp"
	"sort"
)

/*
Find: 从有序数组中查找对应元素

	@params: data 待查找数组
	@params: target 目标元素
	@return: idx 目标元素下标
	@return: exist 是否存在
*/
func Find[T cmp.Ordered](data []T, target T) (idx int, found bool) {
	idx = sort.Search(len(data), func(i int) bool {
		return data[i] >= target
	})

	// 检查是否找到目标元素
	if idx < len(data) && data[idx] == target {
		return idx, true
	}
	return -1, false
}
