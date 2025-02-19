package gray

import (
	"context"
)

// FeatureConfig AB实验配置
type FeatureConfig struct {
	Enable bool           `json:"enable"`
	Rule   []*TrafficRule `json:"rule"` // 分流策略, 影响分组逻辑
}

// Format 格式化配置
func (e *FeatureConfig) Format() {
	for _, rule := range e.Rule {
		rule.Format()
	}
}

// Validate 校验配置
func (e *FeatureConfig) Validate() error {
	var err error
	for _, rule := range e.Rule {
		if err = rule.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// Group 确定分组
func (e *FeatureConfig) Group(ctx context.Context) TrafficGroup {
	for _, rule := range e.Rule {
		if rule.Group(ctx) == TrafficGroupB {
			return TrafficGroupB
		}
	}
	return TrafficGroupA
}
