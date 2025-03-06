package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/everfir/go-helpers/consts"
	"github.com/everfir/go-helpers/define"
	gray_util "github.com/everfir/go-helpers/gray"
	"github.com/everfir/go-helpers/internal/structs/gray"
	"github.com/everfir/logger-go"
	"github.com/everfir/logger-go/structs/log_level"
)

func main() {
	logger.Init(logger.WithLevel(log_level.DebugLevel))
	// 创建一个 GrayConfig 实例并初始化数据
	grayConfig := gray.GrayConfig{
		"business1": {
			Feature: map[string]*gray.FeatureConfig{
				"feature1": {
					Enable: true,
					Rule: []*gray.TrafficRule{
						{
							Enable:      true,
							TargetGroup: "b",
							Rate:        1,
							Targets:     map[string][]string{"platform": []string{"ios"}, "device": []string{"phone"}},
							WhiteList:   []string{"101", "201", "301"},
							BlackList:   []string{"1001", "2001", "3001"},
							Expresion:   "user.ctime >= 1741017900",
							TrafficRate: 0.5,
						},
						// {
						// 	Enable:      true,
						// 	TargetGroup: "c",
						// 	Rate:        0.5,
						// 	Targets: map[string][]string{
						// 		"platform": []string{"ios"},
						// 		"device":   []string{"phone"},
						// 	},
						// 	TrafficRate: 0.5,
						// },
					},
				},
			},
		},
	}

	rand.Seed(time.Now().UnixNano())

	// 统计每个分组的比例
	var mapGroup = make(map[consts.TrafficGroup]int)

	// 构造若干个请求
	var cnt int = 5000
	for i := 0; i < cnt; i++ {
		var ctx context.Context = context.TODO()
		ctx = context.WithValue(ctx, consts.AccountInfoKey, &define.AccountInfo{
			AccountId: uint64(i),
		})
		ctx = context.WithValue(ctx, consts.BusinessKey, "business1")
		ctx = context.WithValue(ctx, consts.PlatformKey, consts.DP_IOS)
		ctx = context.WithValue(ctx, consts.DeviceKey, consts.Dev_Phone)

		// 预期结果≈1:1

		// 开启之后，预期结果≈ 3:1
		if rand.Float64() < 0.5 {
			ctx = context.WithValue(ctx, consts.PlatformKey, consts.DP_Android)
		}

		// 开启之后，预期结果≈ 2:1
		// if rand.Float64() < 0.3 {
		// 	ctx = context.WithValue(ctx, consts.DeviceKey, consts.Dev_PC)
		// }

		group := gray_util.ExperimentGroup(ctx, "feature1", &grayConfig)
		mapGroup[group]++
	}

	for group, count := range mapGroup {
		fmt.Printf("分组 %v 的比例为 %v\n", group, float64(count)/float64(cnt))
	}
}
