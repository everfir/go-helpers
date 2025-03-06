package gray

import (
	"context"

	"github.com/everfir/go-helpers/consts"
	"github.com/everfir/logger-go"
	"github.com/everfir/logger-go/structs/field"
)

type GrayConfig map[string]Gray // key: business

func (gc *GrayConfig) Format() {
	for _, gray := range *gc {
		gray.Format()
	}
}

func (gc *GrayConfig) Validate() error {
	for _, gray := range *gc {
		if err := gray.Validate(); err != nil {
			return err
		}
	}
	return nil
}

type Gray struct {
	Feature map[string]*FeatureConfig `json:"feature"`
}

// Format: 格式化配置
func (g *Gray) Format() {
	for _, config := range g.Feature {
		config.Format()
	}
}

// Validate: 校验配置
func (g *Gray) Validate() error {
	for _, config := range g.Feature {
		if err := config.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (g *Gray) Group(ctx context.Context, feature string) consts.TrafficGroup {
	if _, exist := g.Feature[feature]; !exist {
		return consts.TrafficGroup_A
	}

	return g.Feature[feature].Group(ctx)
}

// Experimental 判断某个功能是否属于某个灰度组
//
// 该方法的逻辑如下：
// 1. 如果功能未配置（Feature 未在 Gray 结构体中定义），认为该功能默认启用（即稳定分支），返回 TrafficGroup_A。
// 2. 如果功能已配置但未启用（Enable 字段为 false），返回 TrafficGroup_A。
// 3. 根据功能的实验分组进行判断：
//   - 如果分组未知（TrafficGroup_Unknow），返回 TrafficGroup_A，并记录警告日志。
//   - 如果分组为 B（TrafficGroup_B），返回 TrafficGroup_B（表示该功能对该分组开放）。
//   - 其他情况返回 TrafficGroup_A（表示该功能对该分组未开放）。
//
// 参数：
//   - ctx: 上下文，用于日志记录和获取用户信息。
//   - feature: 功能名称。
//
// 返回值：
//   - consts.TrafficGroup: 表示该功能的实验分组。
//
// 注意事项：
// 1. 该方法适用于灰度发布场景，通过动态调整实验分组控制功能的开放。
// 2. 功能配置由 Gray 结构体的 Feature 字段管理。
// 3. 分组逻辑由 FeatureConfig.Group 方法确定。
// 4. 确保上下文中包含用户信息，以便正确判断功能的可用性。
//
// 示例：
//
//	if gray.Experimental(ctx, "new_feature") == consts.TrafficGroup_B {
//	    // 执行新功能逻辑
//	} else {
//
//	    // 执行旧逻辑
//	}
func (g Gray) Experimental(ctx context.Context, feature string) consts.TrafficGroup {
	// 检查功能是否已配置
	var exist bool
	var config *FeatureConfig
	config, exist = g.Feature[feature]
	if !exist {
		// 功能未配置，返回稳定分支
		return consts.TrafficGroup_A
	}
	if !config.Enable {
		// 功能已配置但未启用，返回稳定分支
		return consts.TrafficGroup_A
	}

	// 根据用户确定分组
	group := config.Group(ctx)
	if group == consts.TrafficGroup_Unknow {
		// 记录未知分组的警告日志
		logger.Warn(
			ctx,
			"[go-helper] unknown experiment group",
			field.String("group", string(group)),
		)
		return consts.TrafficGroup_A
	}
	// 返回确定的分组
	return group
}
