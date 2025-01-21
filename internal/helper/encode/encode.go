package encode

import "github.com/zeebo/xxh3"

func HashString(str string) uint64 {
	return xxh3.HashString(str)
}
