package encode

import "github.com/zeebo/xxh3"

// HashString 获取xxh3生成的哈希结果
func HashString(str string) uint64 {
	return xxh3.HashString(str)
}
