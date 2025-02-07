package structs

type CheckTokenReq struct {
	Token string `json:"token,omitempty"`
}

type CheckTokenResp struct {
	AccountInfo *AccountInfo `json:"account_info,omitempty"`
	Valid       bool         `json:"valid,omitempty"`
	ErrCode     uint32       `json:"err_code,omitempty"`
	ErrMsg      string       `json:"err_msg,omitempty"`
}

type AccountInfo struct {
	AccountId   uint64 `json:"account_id,omitempty"`
	Role        uint8  `json:"role,omitempty"`
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	Nickname    string `json:"nickname,omitempty"`
	Avatar      string `json:"avatar,omitempty"`
	PhoneNum    string `json:"phone_num,omitempty"`
	Email       string `json:"email,omitempty"`
	Source      uint8  `json:"source,omitempty"`
	Extra       string `json:"extra,omitempty"`
	VipExpireTs uint32 `json:"vip_expire_timestamp,omitempty"`
	Ctime       uint32 `json:"ctime"`
}

type AccountServiceConfig struct {
	Url        string
	GuestToken string
}
