package structs

type AccountInfo struct {
	AccountId     uint64   `json:"account_id"`           // 账户ID，唯一标识符
	Role          uint8    `json:"role"`                 // 用户角色，用于权限控制
	Channel       string   `json:"channel,omitempty"`    // 注册渠道，标识用户来源
	Platform      string   `json:"platform"`             // 用户使用的平台（如iOS、Android）
	Username      string   `json:"username,omitempty"`   // 用户名，用于登录
	Password      string   `json:"password,omitempty"`   // 用户密码（通常加密存储）
	Nickname      string   `json:"nickname,omitempty"`   // 用户昵称，显示名称
	Avatar        string   `json:"avatar,omitempty"`     // 用户头像URL
	PhoneNum      string   `json:"phone_num,omitempty"`  // 用户手机号码
	Email         string   `json:"email,omitempty"`      // 用户邮箱地址
	Source        uint8    `json:"source,omitempty"`     // 用户来源（如APP、Web）
	Extra         string   `json:"extra,omitempty"`      // 额外信息，通常为JSON字符串
	VipExpireTime uint32   `json:"vip_expire_timestamp"` // VIP过期时间戳
	Ctime         uint32   `json:"ctime"`                // 账户创建时间戳
	TemplateIDs   []string `json:"template_ids"`         // 用户关联的模板ID列表
	Business      string   `json:"business"`             // 业务标识，用于区分不同业务线
	WechatUnionId string   `json:"wechat_union_id"`      // 微信UnionID，用于微信生态用户标识
}

func (info *AccountInfo) Validate() bool {
	return info.AccountId != 0
}

func (info *AccountInfo) Param() map[string]interface{} {
	return nil
}
