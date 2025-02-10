package gray

import (
	"context"
	"sync"

	"github.com/everfir/go-helpers/define"
	"github.com/everfir/go-helpers/env"
	"github.com/everfir/go-helpers/internal/helper/nacos"
	"github.com/everfir/go-helpers/internal/structs/gray"
)

var getGrayConfig func() *define.Config[gray.GrayConfig] = sync.OnceValue(func() *define.Config[gray.GrayConfig] {
	config, err := nacos.GetConfigFromNacosAndConfigOnChange[gray.GrayConfig](nacos.GetNacosClient(), "gray.json")
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

	// 业务没有对应的配置，认为此业务是稳定的业务，直接返回true
	if _, exist := config[business]; !exist {
		return true
	}

	return getGrayConfig().Get()[business].Experimental(ctx, feature)
}

// // Experimental 判断某个功能是否处于实验状态。(客户端用)
// //
// // 如果当前业务未找到对应的灰度配置，则默认认为该业务是稳定的，返回 true。
// // 否则，根据业务的灰度配置判断功能是否可用。
// //
// // 参数：
// //   - ctx: 上下文信息，用于获取当前业务环境。
// //   - feature: 需要判断的功能标识符。
// //
// // 注意事项：
// // 1.底层代码会从上下文中获取用户信息，请确保拥有此信息，否则可能导致该功能不可用
// //
// // 返回值：
// //   - bool: 如果该功能处于实验状态，则返回 true；否则返回 false。
// func ExperimentGroup(ctx context.Context, feature string) gray.TrafficGroup {
// 	business := env.Business(ctx)
// 	if business == "" {
// 		return gray.TrafficGroupA
// 	}

// 	config := getGrayConfig().Get()
// 	if _, exist := config[business]; !exist {
// 		return gray.TrafficGroupA
// 	}

// 	if feature == "" {
// 		return gray.TrafficGroupA
// 	}

// 	return config[business].Group(feature)
// }
