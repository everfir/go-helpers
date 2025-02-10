package gray

import (
	"context"

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
func (g Gray) Validate() error {
	for _, config := range g.Feature {
		if err := config.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (g Gray) Group(ctx context.Context, feature string) TrafficGroup {
	if _, exist := g.Feature[feature]; !exist {
		return TrafficGroupA
	}

	return g.Feature[feature].Group(ctx)
}

// Experimental 判断某个功能是否处于实验阶段（灰度发布）
//
// 该方法的逻辑如下：
// 1. 如果功能未配置（Feature 未在 Gray 结构体中定义），认为该功能默认启用（即稳定分支），返回 true。
// 2. 如果功能已配置但未启用（Enable 字段为 false），返回 false。
// 3. 根据功能的实验分组进行判断：
//   - 如果分组未知（TrafficGroupUnKnow），返回 true，并记录警告日志。
//   - 如果分组为 B（TrafficGroupB），返回 true（表示该功能对该分组开放）。
//   - 其他情况返回 false（表示该功能对该分组未开放）。
//
// 参数：
//   - ctx: 上下文，用于日志记录。
//   - feature: 功能名称。
//
// 返回值：
//   - bool: 是否启用该功能。
//
// 注意事项：
// 1. 该方法适用于灰度发布场景，通过动态调整实验分组控制功能的开放。
// 2. 功能配置由 Gray 结构体的 Feature 字段管理。
// 3. 分组逻辑由 FeatureConfig.Group 方法确定。
// 4. 底层代码会从上下文中获取用户信息，请确保拥有此信息，否则可能导致该功能不可用
//
// 示例：
//
//	if gray.Experimental(ctx, "new_feature") {
//	    // 执行新功能逻辑
//	} else {
//	    // 执行旧逻辑
//	}
func (g Gray) Experimental(ctx context.Context, feature string) bool {

	// 如果Feature没有配置，认为是稳定可靠的分支，返回true
	var exist bool
	var config *FeatureConfig
	config, exist = g.Feature[feature]
	if !exist {
		return true
	}
	if !config.Enable {
		return false
	}

	// 根据用户确定分组
	group := config.Group(ctx)
	if group == TrafficGroupUnKnow {
		logger.Warn(
			ctx,
			"[go-helper] unknown experiment group",
			field.String("group", string(group)),
		)
		return true

	}
	return group == TrafficGroupB
}
