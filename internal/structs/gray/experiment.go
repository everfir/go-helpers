package gray

import (
	"context"

	"github.com/everfir/go-helpers/consts"
)

// FeatureConfig: AB实验配置
type FeatureConfig struct {
	Enable bool           `json:"enable"`
	Rule   []*TrafficRule `json:"rule"` // 分流策略, 影响分组逻辑
}

// Format: 格式化配置
func (e *FeatureConfig) Format() {
	for _, rule := range e.Rule {
		rule.Format()
	}
}

// Validate: 校验配置
func (e *FeatureConfig) Validate() error {
	var err error
	for _, rule := range e.Rule {
		if err = rule.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// Group: 根据分流规则确定分组
func (e *FeatureConfig) Group(ctx context.Context) consts.TrafficGroup {
	for _, rule := range e.Rule {
		// 该分流规则已经关闭，跳过
		if !rule.Enable {
			continue
		}

		// 根据分流规则确定分组
		if rule.Group(ctx) {
			return consts.NewTrafficGroupFromString(rule.TargetGroup)
		}
	}

	// 没有匹配到任何分流规则，返回默认分组
	return consts.TrafficGroup_A
}
