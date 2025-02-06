package gray

import (
	"context"
	"sync"

	"github.com/everfir/go-helpers/define"
	"github.com/everfir/go-helpers/env"
	"github.com/everfir/go-helpers/internal/helper/nacos"
	"github.com/everfir/go-helpers/internal/structs"
)

var getGrayConfig func() *define.Config[structs.GrayConfig] = sync.OnceValue(func() *define.Config[structs.GrayConfig] {
	config, err := nacos.GetConfigFromNacosAndConfigOnChange[structs.GrayConfig](nacos.GetNacosClient(), "gray.json")
	if err != nil {
		panic(err.Error())
	}

	return config
})

// Gray: 判断Feature是否可以被访问(结合AB实验 && Feature灰度来判断)
func Gray(ctx context.Context, feature string) bool {
	business := env.Business(ctx)
	if business == "" {
		return false
	}

	config := getGrayConfig().Get()

	// 业务没有对应的配置，认为此业务是稳定的业务，直接返回true
	if _, exist := config[business]; !exist {
		return true
	}

	return getGrayConfig().Get()[business].Enable(ctx, feature)
}

func ExperimentGroup(ctx context.Context, routerKey string) structs.TrafficGroup {
	business := env.Business(ctx)
	if business == "" {
		return structs.TrafficGroupA
	}

	config := getGrayConfig().Get()
	if _, exist := config[business]; !exist {
		return structs.TrafficGroupA
	}

	if routerKey == "" {
		return structs.TrafficGroupA
	}

	return config[business].Group(routerKey)
}
