package gray

import (
	"context"
	"fmt"
	"sort"

	"github.com/everfir/go-helpers/consts"
	"github.com/everfir/go-helpers/define"
	"github.com/everfir/go-helpers/env"
	"github.com/everfir/go-helpers/internal/helper/encode"
	"github.com/everfir/go-helpers/internal/helper/slice"
	"github.com/everfir/logger-go"
	"github.com/everfir/logger-go/structs/field"
	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

// TrafficRule: 分流策略
type TrafficRule struct {
	Enable bool `json:"enable"` // 是否启用

	Rate      float64             `json:"rate"`      // 首次分流比例, 从所有流量中，获取部分流量，用于判断余下的条件
	Expresion string              `json:"expresion"` // 表达式，在Rule模式下生效
	Targets   map[string][]string `json:"targets"`   // 匹配目标
	WhiteList []string            `json:"whitelist"` // 白名单
	BlackList []string            `json:"blacklist"` // 黑名单

	TrafficRate float64 `json:"traffic_rate"` // 分流比例, 满足条件后，分流到指定组的流量比例
	TargetGroup string  `json:"target_group"` // 所属分流组

	expresionProgram *vm.Program
}

// Format: 格式化配置
func (rule *TrafficRule) Format() {
	for _, targets := range rule.Targets {
		sort.Strings(targets)
	}

	sort.Strings(rule.WhiteList)
	sort.Strings(rule.BlackList)
}

// Validate 校验 TrafficRule 的各项配置是否有效。
//
// 该方法执行以下检查：
// 1. 如果规则未启用，则跳过验证。
// 2. 检查 TargetGroup 是否在有效范围内（应为 b-z）。
// 3. 检查 Rate 是否在有效范围内（0 到 1 之间）。
// 4. 检查 TrafficRate 是否在有效范围内（0 到 1 之间）。
// 5. 如果 Expresion 不为空，尝试编译表达式并检查是否有效。
//
// 返回值：
//   - 如果所有检查都通过，返回 nil；
//   - 如果有任何检查失败，返回相应的错误信息。
func (rule *TrafficRule) Validate() error {
	// 如果规则未启用，则跳过验证
	if !rule.Enable {
		return nil
	}

	// 检查 TargetGroup 是否在有效范围内
	if rule.TargetGroup == "" || rule.TargetGroup == "a" {
		return fmt.Errorf("invalid rule.TargetGroup[%s] should be in [b-z]", rule.TargetGroup)
	}

	// 检查 Rate 是否在有效范围内（0 到 1 之间）
	if rule.Rate < 0 || rule.Rate > 1 {
		return fmt.Errorf("invalid rule.Rate[%v] should be in [0, 1]", rule.Rate)
	}

	// 检查 TrafficRate 是否在有效范围内（0 到 1 之间）
	if rule.TrafficRate < 0 || rule.TrafficRate > 1 {
		return fmt.Errorf("invalid rule.TrafficRate[%v] should be in [0, 1]", rule.TrafficRate)
	}

	// 预编译表达式
	if rule.Expresion != "" {
		var err error
		rule.expresionProgram, err = expr.Compile(rule.Expresion, expr.AsBool())
		if err != nil {
			return fmt.Errorf("compile rule.Expresion[%s] failed: %w", rule.Expresion, err)
		}
	}

	// 所有检查都通过
	return nil
}

// Group 根据 TrafficRule 的规则对用户进行分组。
//
// 该函数会根据 TrafficRule 的 Mode 字段进行不同的分组逻辑：
//   - 如果用户在白名单中，则直接返回 true，表示匹配。
//   - 如果用户在黑名单中，则返回 false，表示不匹配。
//   - 先进行首次分流，如果用户的哈希值不符合设定的比例，则返回 false，表示不匹配。
//   - 检查设备信息，如果设备不在目标设备列表中，则返回 false，表示不匹配。
//   - 检查平台信息，如果平台不在目标平台列表中，则返回 false，表示不匹配。
//   - 检查应用类型，如果应用类型不在目标应用类型列表中，则返回 false，表示不匹配。
//   - 如果定义了表达式，则运行预编译的表达式并根据结果返回匹配状态。
//   - 最后进行二次分流，如果用户的哈希值不符合设定的流量比例，则返回 false，表示不匹配。
//   - 如果所有检查都通过，则返回 true，表示匹配。
//
// 参数：
//   - ctx: 上下文对象，用于传递请求上下文信息。
//
// 返回值：
//   - bool: 如果用户匹配该规则则返回 true，否则返回 false。
//
// 示例：
//
//	rule := &TrafficRule{
//	    Mode:      TrafficModeRate,
//	    Rate:      0.5,
//	    WhiteList: []string{"123", "456"},
//	}
//	match := rule.Group(context.Background())
//	fmt.Println(match) // true 或 false
func (rule *TrafficRule) Group(ctx context.Context) (match bool) {
	// 需要参与判断的数据
	var device consts.TDevice = env.Device(ctx)               // 获取设备信息
	var appType consts.TAppType = env.AppType(ctx)            // 获取应用类型
	var platform consts.TDevicePlatform = env.Platform(ctx)   // 获取平台信息
	var accountInfo define.AccountInfo = env.AccountInfo(ctx) // 获取用户账户信息

	// 记录当前分组信息的调试日志
	logger.Debug(
		ctx,
		"group info",
		field.Any("device", device),
		field.Any("appType", appType),
		field.Any("platform", platform),
		field.Any("accountInfo", accountInfo),
	)

	// 检查白名单
	accountId := fmt.Sprintf("%d", accountInfo.AccountId) // 将账户 ID 转换为字符串
	if _, exist := slice.Find(rule.WhiteList, accountId); exist {
		// 如果账户在白名单中，返回 true，表示匹配
		return true
	}

	// 检查黑名单
	if _, exist := slice.Find(rule.BlackList, accountId); exist {
		// 如果账户在黑名单中，返回 false，表示不匹配
		return false
	}

	// 先进行首次分流
	if !rateMatch(accountId, rule.Rate) {
		// 如果用户的哈希值不符合设定的比例，返回 false，表示不匹配
		return false
	}

	// 检查设备
	if len(rule.Targets["device"]) > 0 {
		if _, exist := slice.Find(rule.Targets["device"], device.String()); !exist {
			// 如果设备不在目标设备列表中，返回 false，表示不匹配
			return false
		}
	}

	// 检查平台
	if len(rule.Targets["platform"]) > 0 {
		if _, exist := slice.Find(rule.Targets["platform"], platform.String()); !exist {
			// 如果平台不在目标平台列表中，返回 false，表示不匹配
			return false
		}
	}

	// 检查应用类型
	if len(rule.Targets["app_type"]) > 0 {
		if _, exist := slice.Find(rule.Targets["app_type"], appType.String()); !exist {
			// 如果应用类型不在目标应用类型列表中，返回 false，表示不匹配
			return false
		}
	}

	// 检查表达式
	if rule.Expresion != "" && rule.expresionProgram != nil {
		param := makeParam(ctx, &accountInfo)                // 创建表达式参数
		match, err := expr.Run(rule.expresionProgram, param) // 运行表达式
		if err != nil {
			// 如果表达式运行失败，记录警告日志并返回 false，表示不匹配
			logger.Warn(
				ctx,
				"expr.Run failed",
				field.String("error", err.Error()),
				field.String("rule.Expresion", rule.Expresion),
				field.Any("param", param),
			)
			return false
		}
		if !match.(bool) {
			// 如果表达式结果为 false，返回 false，表示不匹配
			return false
		}
	}

	// 二次分流
	if !rateMatch(accountId, rule.TrafficRate) {
		// 如果用户的哈希值不符合设定的流量比例，返回 false，表示不匹配
		return false
	}

	// 所有检查都通过，返回 true，表示匹配
	return true
}

func rateMatch(accountId string, rate float64) bool {
	hash := encode.HashString(accountId)
	bucket := hash % 1000
	threshold := uint64(rate * float64(1000))
	return bucket < threshold
}

func makeParam(ctx context.Context, accountInfo *define.AccountInfo) (ret map[string]interface{}) {
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
