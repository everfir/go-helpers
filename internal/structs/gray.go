package structs

import (
	"sort"
)

type Gray map[string]map[string]GrayConfig

func (g Gray) Enable(business, feature, user string) bool {

	config, ok := g[business]
	if !ok {
		return false
	}

	featConfig, ok := config[feature]
	if !ok {
		return false
	}

	return featConfig.White(user)
}

func (g *Gray) Format() {
	for _, config := range *g {
		for _, feat := range config {
			feat.Format()
		}
	}
}

type GrayConfig struct {
	Enable    bool     `json:"enable"`
	WhiteList []string `json:"whitelist"`
}

func (c *GrayConfig) Format() {
	if !c.Enable {
		return
	}

	sort.Strings(c.WhiteList)
}

func (c *GrayConfig) White(user string) bool {
	if !c.Enable {
		return false
	}

	// 空白名单列表，表示对所有人开放
	if len(c.WhiteList) == 0 {
		return true
	}

	idx := sort.SearchStrings(c.WhiteList, user)
	return idx < len(c.WhiteList) && c.WhiteList[idx] == user
}
