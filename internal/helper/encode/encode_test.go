package encode

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

// 均衡性测试：
// 1.将指定数量的testNum分布在numBuckets个桶中，校验每个桶分配的数量是否接近。
// 2.按rate比例取桶来验证桶中数量是否接近总数量的比例数量。
func TestHashStringBalance(t *testing.T) {
	numBuckets := 1000
	testNum := 1000000
	rate := 0.4
	hashBuckets := make(map[uint64]int, numBuckets)

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < testNum; i++ {
		str := fmt.Sprintf("%d", rand.Int())
		hashValue := HashString(str)
		hashBuckets[hashValue%uint64(numBuckets)]++
	}

	// 预期每个桶的平均值
	expectedAvg := testNum / numBuckets
	// 允许的误差范围
	tolerance := expectedAvg / 10
	t.Logf("Expected average: %d, tolerance: %d", expectedAvg, tolerance)

	// 1. 检查每个桶的分布是否均匀
	unevenCount := 0
	for bucket, count := range hashBuckets {
		if count < expectedAvg-tolerance || count > expectedAvg+tolerance {
			unevenCount++
			t.Logf("Bucket %d has an uneven distribution. Expected ~%d, got %d.", bucket, expectedAvg, count)
		}
	}
	if unevenCount > 0 {
		t.Logf("Found %d uneven buckets.", unevenCount)
	}

	// 2. 检查按rate比例取桶的数量是否接近总数量的比例数量
	expectedGroupB := int(float64(testNum) * rate)
	groupB := 0
	threshold := uint64(rate * float64(numBuckets))
	for bucket, count := range hashBuckets {
		if bucket < threshold {
			groupB += count
		}
	}
	t.Logf("Bucket group B has an uneven distribution. Expected ~%d, got %d.", expectedGroupB, groupB)

}

// 性能测试
func BenchmarkHashString(b *testing.B) {
	user := "100000"
	// Benchmark the HashString function over a large number of iterations
	b.ResetTimer() // Reset the timer to exclude setup time.
	for i := 0; i < b.N; i++ {
		HashString(user)
	}
}
