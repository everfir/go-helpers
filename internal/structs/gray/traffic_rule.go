package gray

import (
	"context"
	"fmt"
	"sort"

	"github.com/everfir/go-helpers/env"
	"github.com/everfir/go-helpers/internal/helper/encode"
	"github.com/everfir/go-helpers/internal/helper/slice"
	"github.com/everfir/go-helpers/internal/structs"
	"github.com/everfir/logger-go"
	"github.com/everfir/logger-go/structs/field"
	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

type TrafficMode uint8

const (
	TrafficModeRate     TrafficMode = iota // 按比例分流(用户ID)
	TrafficModeRule                        // 按规则分类(表达式匹配)
	TrafficModeDevice                      // 按客户端设备分流 phone/pc/ipad
	TrafficModeVersion                     // 按版本分流
	TrafficModePlatform                    // 按平台分流 android/ios/linux/windows/macos/ipados

	TrafficModeUnknow = iota + 1 // 按平台分流 android/ios/linux/windows/macos/ipados
)

type TrafficGroup string

const (
	TrafficGroupA      TrafficGroup = "A"
	TrafficGroupB      TrafficGroup = "B"
	TrafficGroupUnKnow TrafficGroup = ""
)

// TrafficRule: 分流策略
type TrafficRule struct {
	Mode      TrafficMode `json:"mode"`      // 分流模式
	Rate      float64     `json:"rate"`      // 分流比例, 仅在Rate模式下生效
	Expresion string      `json:"expresion"` // 表达式，在Rule模式下生效
	Targets   []string    `json:"targets"`   // 匹配目标，在Device/Version/Platform
	WhiteList []string    `json:"whitelist"` // 白名单, 全局有效

	expresionProgram *vm.Program
}

// Format: 格式化配置
func (rule *TrafficRule) Format() {
	sort.Strings(rule.Targets)
	sort.Strings(rule.WhiteList)
}

// Validate 验证 TrafficRule 对象的有效性。
//
// 该函数会根据 TrafficRule 的 Mode 字段进行不同的验证逻辑：
//   - 如果 Mode 是 TrafficModeRate，则检查 Rate 是否在有效范围内（0 到 1 之间）。
//   - 如果 Mode 是 TrafficModeRule，则尝试编译表达式，如果编译失败则返回错误。
//   - 如果 Mode 是 TrafficModeDevice、TrafficModeVersion 或 TrafficModePlatform，则检查 Targets 是否为空。
//   - 如果 Mode 不在有效范围内（小于 TrafficModeRate 或大于等于 TrafficModeUnknow），则返回错误。
//
// 返回值：
//   - 如果 TrafficRule 有效，则返回 nil。
//   - 如果 TrafficRule 无效，则返回具体的错误信息。
//
// 示例：
//
//	rule := &TrafficRule{
//	    Mode:  TrafficModeRate,
//	    Rate:  0.5,
//	}
//	err := rule.Validate()
//	if err != nil {
//	    log.Fatalf("Validation failed: %v", err)
//	}
//
//	rule := &TrafficRule{
//	    Mode:      TrafficModeRule,
//	    Expresion: "user.age > 18",
//	}
//	err := rule.Validate()
//	if err != nil {
//	    log.Fatalf("Validation failed: %v", err)
//	}
//
//	rule := &TrafficRule{
//	    Mode:    TrafficModeDevice,
//	    Targets: []string{"phone", "tablet"},
//	}
//	err := rule.Validate()
//	if err != nil {
//	    log.Fatalf("Validation failed: %v", err)
//	}
func (rule *TrafficRule) Validate() error {
	// 检查 TrafficMode 是否在有效范围内
	// 如果 rule.Mode 小于 TrafficModeRate 或大于等于 TrafficModeUnknow，则返回错误
	if rule.Mode < TrafficModeRate || rule.Mode >= TrafficModeUnknow {
		return fmt.Errorf("Invalid rule.Mode[%v] should be in [%d, %d)", rule.Mode, TrafficModeRate, TrafficModeUnknow)
	}

	// 根据 TrafficMode 进行具体验证
	switch rule.Mode {
	case TrafficModeRate:
		// 如果 TrafficMode 是 TrafficModeRate，检查 rate 是否在有效范围内（0 到 1 之间）
		if rule.Rate < 0 || rule.Rate > 1 {
			return fmt.Errorf("Invalid rule.Rate[%v] should be in [0, 1]", rule.Rate)
		}

	case TrafficModeRule:
		// 如果 TrafficMode 是 TrafficModeRule，编译表达式并检查是否有效
		var err error
		rule.expresionProgram, err = expr.Compile(rule.Expresion, expr.AsBool())
		if err != nil {
			return fmt.Errorf("Compile rule.Expresion[%s] failed: %w", rule.Expresion, err)
		}

	case TrafficModeDevice, TrafficModeVersion, TrafficModePlatform:
		// 如果是 TrafficModeDevice、TrafficModeVersion 或 TrafficModePlatform，检查 Targets 是否为空
		if len(rule.Targets) == 0 {
			return fmt.Errorf("Invalid rule.Targets: should not be empty")
		}

	default:
		// 如果 TrafficMode 是未知值，返回错误
		return fmt.Errorf("Unknown TrafficMode: %v", rule.Mode)
	}

	// 所有检查都通过
	return nil
}

// Group 根据 TrafficRule 的规则对用户进行分组。
//
// 该函数会根据 TrafficRule 的 Mode 字段进行不同的分组逻辑：
//   - 如果用户在白名单中，则直接分组为 TrafficGroupB。
//   - 如果 Mode 是 TrafficModeRate，则根据用户 ID 的哈希值和设定的比例进行分组。
//   - 如果 Mode 是 TrafficModeRule，则运行预编译的表达式并根据结果分组。
//   - 如果 Mode 是 TrafficModeDevice、TrafficModeVersion 或 TrafficModePlatform，则根据设备、版本或平台信息进行分组。
//   - 如果 Mode 无效，则记录警告日志并返回默认分组 TrafficGroupA。
//
// 参数：
//   - ctx: 上下文对象，用于传递请求上下文信息。
//
// 返回值：
//   - TrafficGroup: 分组结果，默认为 TrafficGroupA。
//
// 示例：
//
//	rule := &TrafficRule{
//	    Mode:      TrafficModeRate,
//	    Rate:      0.5,
//	    WhiteList: []string{"123", "456"},
//	}
//	group := rule.Group(context.Background())
//	fmt.Println(group) // TrafficGroupA 或 TrafficGroupB
func (rule *TrafficRule) Group(ctx context.Context) (group TrafficGroup) {
	group = TrafficGroupA

	var accountInfo structs.AccountInfo = env.AccountInfo(ctx)

	// 如果在白名单中，直接分组为B
	if _, exist := slice.Find(rule.WhiteList, fmt.Sprintf("%d", accountInfo.AccountId)); exist {
		return TrafficGroupB
	}
	switch rule.Mode {
	case TrafficModeRate:
		if rule.rateHit(accountInfo.AccountId) {
			group = TrafficGroupB
		}

	case TrafficModeRule:
		param := makeParam(ctx, &accountInfo)
		param["SliceHas"] = contains
		val, err := expr.Run(rule.expresionProgram, param)
		if err != nil {
			logger.Warn(
				context.TODO(),
				"run expresion failed",
				field.String("err", err.Error()),
				field.String("expresion", rule.Expresion),
				field.Any("param", param),
			)
			break
		}

		if val.(bool) && rule.rateHit(accountInfo.AccountId) {
			group = TrafficGroupB
		}

	case TrafficModeDevice:
		device := env.Device(ctx)
		if _, exist := slice.Find[string](rule.Targets, string(device)); exist && rule.rateHit(accountInfo.AccountId) {
			group = TrafficGroupB
		}

	case TrafficModeVersion:
		version := env.Version(ctx)
		if _, exist := slice.Find[string](rule.Targets, string(version)); exist && rule.rateHit(accountInfo.AccountId) {
			group = TrafficGroupB
		}
	case TrafficModePlatform:
		platform := env.Platform(ctx)
		if _, exist := slice.Find[string](rule.Targets, string(platform)); exist && rule.rateHit(accountInfo.AccountId) {
			group = TrafficGroupB
		}
	default:
		logger.Warn(
			context.TODO(),
			"invalid TrafficMode",
			field.Uint8("mode", uint8(rule.Mode)),
		)
	}

	return group
}

func (rule *TrafficRule) rateHit(accountId uint64) bool {
	if rule.Rate == 0 {
		return false
	}

	hash := encode.HashString(fmt.Sprintf("%d", accountId))
	bucket := hash % 1000
	threshold := uint64(rule.Rate * float64(1000))

	if bucket < threshold {
		return true
	}

	return false
}

func makeParam(ctx context.Context, accountInfo *structs.AccountInfo) (ret map[string]interface{}) {
	templateIds := make([]interface{}, 0, len(accountInfo.TemplateIDs))
	for _, id := range accountInfo.TemplateIDs {
		templateIds = append(templateIds, id)
	}

	ret = make(map[string]interface{})
	m := make(map[string]interface{})
	m["account_id"] = accountInfo.AccountId
	m["role"] = accountInfo.Role
	m["channel"] = accountInfo.Channel
	m["platform"] = accountInfo.Platform
	m["username"] = accountInfo.Username
	m["password"] = accountInfo.Password
	m["nickname"] = accountInfo.Nickname
	m["avatar"] = accountInfo.Avatar
	m["phone_num"] = accountInfo.PhoneNum
	m["email"] = accountInfo.Email
	m["source"] = accountInfo.Source
	m["extra"] = accountInfo.Extra
	m["vip_expire_timestamp"] = accountInfo.VipExpireTime
	m["ctime"] = accountInfo.Ctime
	m["template_ids"] = templateIds
	m["business"] = accountInfo.Business
	m["wechat_union_id"] = accountInfo.WechatUnionId

	d := make(map[string]interface{}, 0)
	d["device"] = env.Device(ctx).String()
	d["platform"] = env.Platform(ctx).String()
	d["version"] = env.Version(ctx)
	d["app_type"] = env.AppType(ctx).String()

	ret["user"] = m
	ret["app"] = d
	return ret
}

// 自定义一个 contains 函数，判断数组是否包含某个元素
func contains(params ...any) bool {
	if len(params) != 2 {
		return false
	}

	// 获取数组和要查找的元素
	arr, ok := params[0].([]any)
	if !ok {
		return false
	}

	value := params[1]

	// 检查数组中是否包含指定元素
	for _, v := range arr {
		if v == value {
			return true
		}
	}

	return false
}
