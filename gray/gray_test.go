package gray

import (
	"context"
	"os"
	"testing"

	"github.com/everfir/go-helpers/consts"
	"github.com/everfir/go-helpers/env"
)

// 初始化测试环境
func init() {
	// 模拟环境变量设置
	_ = os.Setenv("BUSINESS", "test_business")
}

// BenchmarkExperimental 性能测试
func BenchmarkExperimental(b *testing.B) {
	// 准备测试上下文
	ctx := context.WithValue(context.Background(), consts.BusinessKey, "test_business")

	// 测试用例
	testCases := []struct {
		name    string
		feature string
	}{
		{"ExistingFeature", "feature_1"},
		{"NonExistingFeature", "unknown_feature"},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Experimental(ctx, tc.feature)
			}
		})
	}
}

// TestExperimental 功能测试
func TestExperimental(t *testing.T) {
	ctx := context.WithValue(context.Background(), consts.BusinessKey, "test_business")

	tests := []struct {
		name     string
		feature  string
		expected bool
	}{
		{"ExistingEnabledFeature", "feature_test", true},
		{"ExistingDisabledFeature", "feature_2", true},
		{"NonExistingFeature", "unknown_feature", true},
		{"EmptyBusiness", "", false},
	}

	ctx = env.TestSetAccountInfo(ctx, uint64(1))
	ctx = context.WithValue(ctx, consts.BusinessKey, "helper_test")
	ctx = context.WithValue(ctx, consts.PlatformKey, consts.DP_MacOS)
	ctx = context.WithValue(ctx, consts.DeviceKey, consts.Dev_PC)
	ctx = context.WithValue(ctx, consts.VersionKey, "v0.0.1")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 测试空业务场景
			if tt.name == "EmptyBusiness" {
				ctx = context.Background()
			}

			if got := Experimental(ctx, tt.feature); got != tt.expected {
				t.Errorf("Experimental(%v) = %v, want %v", tt.feature, got, tt.expected)
			}
		})
	}
}
