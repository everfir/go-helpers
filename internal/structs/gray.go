package structs

import (
	"context"
	"sort"

	"github.com/everfir/go-helpers/env"
	"github.com/everfir/go-helpers/internal/helper/encode"
	"github.com/everfir/go-helpers/internal/helper/slice"
	"github.com/everfir/logger-go"
	"github.com/everfir/logger-go/structs/field"
)

type TrafficMode uint8

const (
	TrafficRuleRate      TrafficMode = iota // 按比例分流
	TrafficRuleWhiteList                    // 白名单分流
	TrafficRuleFusion                       // 混合使用
)

type TrafficGroup string

const (
	TrafficGroupA      TrafficGroup = "A"
	TrafficGroupB      TrafficGroup = "B"
	TrafficGroupUnKnow TrafficGroup = ""
)

type GrayConfig map[string]Gray

func (gc *GrayConfig) Format() {
	for _, gray := range *gc {
		gray.Format()
	}
}

func (gc *GrayConfig) Validate() bool {
	for _, gray := range *gc {
		if !gray.Validate() {
			return false
		}
	}
	return true
}

type Gray struct {
	Experiment  Experiment      `json:"experiment"`   // AB实验配置
	FeatureGray map[string]bool `json:"feature_gray"` // Feature灰度配置
}

// Format: 格式化配置
func (g *Gray) Format() {
	g.Experiment.Format()
}

// Validate: 校验配置
func (g Gray) Validate() bool {
	return g.Experiment.Validate()
}

func (g Gray) Group(user string) TrafficGroup {
	return g.Experiment.Group(user)
}

/*
Enable: 判断功能是否开启(结合AB实验和Feature灰度来判断)

## AB实验判断

- 如果功能不在可灰度列表中，则认为功能是稳定功能，进入Feature灰度判断
- 如果功能在可灰度列表中，比对用户分组和对应分组中可用的功能列表, 当前用户分组不在可用功能列表中，则返回false，否则进入Feature灰度判断

## Feature灰度判断

- 如果Feature灰度配置为空，则认为功能是稳定功能，返回true
- 如果Feature灰度配置不为空，则根据配置T/F来决定功能是否可用
*/
func (g Gray) Enable(ctx context.Context, feature string) bool {
	// 如果Feature没有在实验中，说明是稳定功能，则查看功能是否可用
	if !g.Experiment.Experimental(feature) {
		return g.AvailableFeature(ctx, feature)
	}

	// 如果Feature在实验中，则比对分组和对应的Feature列表
	group := TrafficGroup(env.ExperimentGroup(ctx))
	switch group {
	case TrafficGroupA, TrafficGroupB:
		return g.Experiment.ExperimentalGroup(group, feature) && g.AvailableFeature(ctx, feature)
	default:
		logger.Warn(ctx, "[go-helper] unknown experiment group", field.String("group", string(group)))
		// return false
	}

	return false
}

// AvailableFeature: 判断功能是否可用
func (g *Gray) AvailableFeature(ctx context.Context, feature string) bool {

	var exist bool
	// 如果没有在灰度列表中，认为是正常功能，正常使用
	_, exist = g.FeatureGray[feature]
	if !exist {
		return true
	}

	return g.FeatureGray[feature]
}

// Experiment: AB实验配置
type Experiment struct {
	Rule        TrafficRule `json:"rule"`         // 分流策略
	FeatureList []string    `json:"feature_list"` // 可AB实验的Feature列表

	ExperimentsA []string `json:"experiments_a"` // A组实验中的Feature列表
	ExperimentsB []string `json:"experiments_b"` // A组实验中的Feature列表
}

// Format: 格式化配置
func (e *Experiment) Format() {
	e.Rule.Format()

	sort.Strings(e.FeatureList)
	sort.Strings(e.ExperimentsA)
	sort.Strings(e.ExperimentsB)
}

// Validate: 校验配置
func (e *Experiment) Validate() bool {
	if !e.Rule.Validate() {
		return false
	}

	featSlices := [][]string{
		e.FeatureList,
		e.ExperimentsA,
		e.ExperimentsB,
	}

	for _, features := range featSlices {
		for _, feature := range features {
			if feature == "" {
				return false
			}
		}
	}

	return true
}

func (e *Experiment) White(user string) bool {
	return e.Rule.White(user)
}

func (e *Experiment) Group(user string) TrafficGroup {
	return e.Rule.Group(user)
}

// Experimental: 判断feature是否在在实验中
func (e *Experiment) Experimental(feature string) bool {
	_, exist := slice.Find[string](e.FeatureList, feature)
	return exist
}

// ExperimentalGroup: 判断feature是否在实验组中
func (e *Experiment) ExperimentalGroup(group TrafficGroup, feature string) bool {
	var exist bool
	if group == TrafficGroupA {
		_, exist = slice.Find[string](e.ExperimentsA, feature)
	} else if group == TrafficGroupB {
		_, exist = slice.Find[string](e.ExperimentsB, feature)
	}

	return exist
}

// TrafficRule: 分流策略
type TrafficRule struct {
	Mode      TrafficMode `json:"mode"`      // 分流模式
	Rate      float64     `json:"rate"`      // 分流比例
	WhiteList []string    `json:"whitelist"` // 白名单
}

// Format: 格式化配置
func (rule *TrafficRule) Format() {
	sort.Strings(rule.WhiteList)
}

// Validate: 校验配置
func (rule *TrafficRule) Validate() bool {
	if rule.Mode < TrafficRuleRate || rule.Mode > TrafficRuleFusion {
		return false
	}

	if rule.Rate < 0 || rule.Rate > 1 {
		return false
	}

	return true
}

func (rule *TrafficRule) White(user string) bool {
	_, exist := slice.Find[string](rule.WhiteList, user)
	return exist
}

func (rule *TrafficRule) Group(user string) (group TrafficGroup) {
	group = TrafficGroupA

	switch rule.Mode {
	case TrafficRuleWhiteList:
		if rule.White(user) {
			group = TrafficGroupB
		}
	case TrafficRuleRate:
		if rule.Rate == 0 {
			break
		}

		hash := encode.HashString(user)
		bucket := hash % 1000
		threshold := uint64(rule.Rate * float64(1000))

		if bucket < threshold {
			group = TrafficGroupB
		}
	case TrafficRuleFusion:
		if rule.White(user) {
			group = TrafficGroupB
		}

		if rule.Rate == 0 {
			break
		}

		hash := encode.HashString(user)
		if (hash % 10) <= uint64(rule.Rate*10) {
			group = TrafficGroupB
		}

	default:
		break
	}

	return group
}
