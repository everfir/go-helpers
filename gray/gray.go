package gray

import (
	"context"
	"sync"

	"github.com/everfir/go-helpers/define/config"
	"github.com/everfir/go-helpers/env"
	"github.com/everfir/go-helpers/internal/helper/nacos"
	"github.com/everfir/go-helpers/internal/structs/gray"
)

var getGrayConfig func() *config.NacosConfig[gray.GrayConfig] = sync.OnceValue(func() *config.NacosConfig[gray.GrayConfig] {
	config, err := nacos.GetConfigAndListen[gray.GrayConfig](nacos.GetNacosClient(), "gray.json")
	if err != nil {
		panic(err.Error())
	}

	return config
})

// Experimental 判断某个功能(feature)在当前业务环境下是否处于实验阶段。
// 如果业务标识为空，则默认返回 false，表示不可用。
// 如果业务没有对应的灰度配置，则认为该业务是稳定业务，默认返回 true。
// 否则，调用具体业务的灰度实验配置进行判断。
func Experimental(ctx context.Context, feature string) bool {
	business := env.Business(ctx)
	if business == "" {
		return false
	}

	config := getGrayConfig().Get()

	// 业务没有对应的配置，认为此业务是稳定的业务，直接返回 false
	if _, exist := config[business]; !exist {
		return false
	}

	return getGrayConfig().Get()[business].Experimental(ctx, feature)
}

// GetAllEnableFeature 获取所有启动状态的feat名称
func GetAllEnableFeature(ctx context.Context) []string {
	business := env.Business(ctx)
	if business == "" {
		return nil
	}

	config := getGrayConfig().Get()

	// 业务没有对应的配置，认为此业务是稳定的业务，直接返回 false
	if _, exist := config[business]; !exist {
		return nil
	}

	ret := make([]string, 0, len(config[business].Feature))
	for feat, rule := range config[business].Feature {
		if rule.Enable {
			ret = append(ret, feat)
		}
	}

	return ret
}
