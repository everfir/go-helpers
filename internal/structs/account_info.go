package structs

type AccountConfig struct {
	UrlEnv map[string]string `json:"url_env"`
}

type AccountInfo struct {
	AccountId     uint64   `json:"account_id,omitempty"`
	Role          uint8    `json:"role"`
	Channel       string   `json:"channel,omitempty"`
	Platform      string   `json:"platform"`
	Username      string   `json:"username,omitempty"`
	Password      string   `json:"password,omitempty"`
	Nickname      string   `json:"nickname,omitempty"`
	Avatar        string   `json:"avatar,omitempty"`
	PhoneNum      string   `json:"phone_num,omitempty"`
	Email         string   `json:"email,omitempty"`
	Source        uint8    `json:"source,omitempty"`
	Extra         string   `json:"extra,omitempty"`
	VipExpireTime uint32   `json:"vip_expire_timestamp"`
	Ctime         uint32   `json:"ctime"`
	TemplateIDs   []string `json:"template_ids"`
	Business      string   `json:"business"`
	WechatUnionId string   `json:"wechat_union_id"`
}

func (info *AccountInfo) Validate() bool {
	return info.AccountId != 0
}
