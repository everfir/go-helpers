package nacos

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	nacos_config "github.com/everfir/go-helpers/define/config"
	"github.com/everfir/go-helpers/env"
	"github.com/everfir/go-helpers/internal/structs"
	internal_config "github.com/everfir/go-helpers/internal/structs/config"
	"github.com/everfir/logger-go"
	"github.com/everfir/logger-go/structs/field"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

var GetNacosClient func() config_client.IConfigClient = sync.OnceValue(func() config_client.IConfigClient {
	namespace := Namespace()
	ipAddr := NacosIp()
	username, passward := AuthInfo()

	cc := constant.ClientConfig{
		NamespaceId:         namespace,
		TimeoutMs:           60000,
		NotLoadCacheAtStart: true,
		LogDir:              "",
		CacheDir:            "",
		LogLevel:            "error",
		Username:            username,
		Password:            passward,
	}

	sc := []constant.ServerConfig{
		{
			IpAddr: ipAddr,
			Port:   8848,
		},
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"clientConfig":  cc,
		"serverConfigs": sc,
	})
	if err != nil {
		panic(fmt.Sprintf("[go-helper] Init nacos client failed: %v ip:%s", err, ipAddr))
	}

	return configClient
})

// GetConfigAndListen 从 Nacos 获取配置并监听配置变更
//
// 参数:
//   - client: Nacos 配置客户端，用于与 Nacos 服务器进行交互
//   - dataId: 配置的唯一标识，用于指定要获取的配置项
//
// 返回值:
//   - config: 包含配置数据的 NacosConfig 对象，
//     该对象的 Data 字段会根据环境存储配置
//   - err: 错误信息，如果获取配置或监听过程中发生错误，则返回相应的错误信息
//
// 使用场景:
//   - 当需要从 Nacos 获取特定配置并实时监听配置变更时，可以使用该函数
//   - 适用于动态配置管理的场景，例如微服务架构中的配置管理
//
// 示例：
//
//	import (
//	    "github.com/everfir/go-helpers/define/config"
//	    "log"
//	)
//
//	var err error
//	var cfg *config.NacosConfig[AppConfig]
//	cfg, err = GetConfigAndListen[AppConfig](client, "app-config")
//	if err != nil {
//	    log.Fatal(err)
//	}
func GetConfigAndListen[T any](client config_client.IConfigClient, dataId string) (config *nacos_config.NacosConfig[T], err error) {
	var conf *internal_config.Config[T]
	conf, err = getConfigAndListen[T](client, dataId, env.Env())
	if err != nil {
		return
	}

	data := make(map[string]*internal_config.Config[T])
	data[env.Env()] = conf
	return nacos_config.NewNacosConfig[T](data), nil
}

// GetConfigAndListenWithGray 从 Nacos 获取配置并监听配置变更，
// 会自动寻找与灰度相关的配置。
//
// 参数:
//   - client: Nacos 配置客户端，用于与 Nacos 服务器进行交互
//   - dataId: 配置的唯一标识，用于指定要获取的配置项
//
// 返回值:
//   - config: 包含所有灰度配置数据的 NacosConfig 对象，
//     该对象的 Data 字段会根据环境存储配置
//   - err: 错误信息，如果获取配置或监听过程中发生错误，则返回相应的错误信息
//
// 使用场景:
//   - 当需要根据灰度配置获取配置时，可以使用该函数
//   - 适用于动态配置管理的场景，例如微服务架构中的灰度发布
//
// 示例：
//
//	import (
//	    "github.com/everfir/go-helpers/define/config"
//	    "log"
//	)
//
//	var err error
//	var cfg *config.NacosConfig[AppConfig]
//	cfg, err = GetConfigAndListenWithGray[AppConfig](client, "app-config")
//	if err != nil {
//	    log.Fatal(err)
//	}
func GetConfigAndListenWithGray[T any](
	client config_client.IConfigClient,
	dataId string,
) (config *nacos_config.NacosConfig[T], err error) {
	logger.Debug(context.Background(), "GetConfigAndListenWithGray", field.String("dataId", dataId))

	// 搜索该 DataId 下所有的 Group 配置
	searchConfigParam := vo.SearchConfigParam{
		Search: "accurate", // 精确搜索
		DataId: dataId,     // 指定要查询的 DataId
		// 这里可以添加其他参数，例如 Namespace，如果需要的话
	}

	// 执行搜索以获取所有相关配置
	searchConfig, err := client.SearchConfig(searchConfigParam)
	if err != nil {
		// 处理错误，记录日志并返回
		return nil, fmt.Errorf("[go-helper] Search config failed, err: %w", err)
	}

	// 将搜索到的配置按 Group 分组
	var group2Config = make(map[string]struct{})
	for _, config := range searchConfig.PageItems {
		group2Config[config.Group] = struct{}{}
	}
	logger.Debug(context.Background(), "group2Config", field.Any("group2Config", group2Config))

	// 获取所有的 DataId 配置
	var data = make(map[string]*internal_config.Config[T])
	for group, _ := range group2Config {
		logger.Debug(context.Background(), "getConfigAndListen", field.String("dataId", dataId), field.String("group", group))
		// 获取配置并监听配置变更
		var conf *internal_config.Config[T]
		conf, err = getConfigAndListen[T](client, dataId, group)
		if err != nil {
			return nil, fmt.Errorf("[go-helper] Get config and listen failed for group %s, err: %w", group, err)
		}
		logger.Info(
			context.Background(),
			"getConfigAndListen",
			field.String("dataId", dataId),
			field.String("group", group),
			field.Any("conf", conf.Get()),
		)

		// 将获取到的配置存储到数据映射中
		data[group] = conf
	}

	// 返回包含所有灰度配置的 NacosConfig 对象
	return nacos_config.NewNacosConfig[T](data), nil
}

// getConfigAndListen 从 Nacos 获取配置并监听配置变更
//
// T: 配置结构体类型，表示要解析的配置数据类型
//
// 参数:
//   - client: Nacos 配置客户端，用于与 Nacos 服务器进行交互
//   - dataId: 配置的唯一标识，用于指定要获取的配置项
//   - group: 配置所属的分组，用于区分不同的配置组
//
// 返回值:
//   - config: 包含配置数据的 Config 对象，
//     该对象的 Data 字段会根据环境存储配置
//   - err: 错误信息，如果获取配置或监听过程中发生错误，则返回相应的错误信息
//
// 使用场景:
//   - 当需要从 Nacos 获取特定配置并实时监听配置变更时，可以使用该函数
//   - 适用于动态配置管理的场景，例如微服务架构中的配置管理
func getConfigAndListen[T any](
	client config_client.IConfigClient,
	dataId string,
	group string,
) (config *internal_config.Config[T], err error) {
	// 从 Nacos 获取配置
	cfg, err := client.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
	logger.Debug(
		context.Background(),
		"getConfigAndListen",
		field.String("dataId", dataId),
		field.String("group", group),
		field.String("cfg", cfg),
	)

	// 如果获取配置时发生错误，返回错误信息
	if err != nil {
		err = fmt.Errorf("[go-helper] Get config from nacos failed, err: %w", err)
		return
	}

	// 创建新的配置对象并解析 JSON
	config = internal_config.NewConfig[T]()
	err = json.Unmarshal([]byte(cfg), config.Data)
	if err != nil {
		return nil, fmt.Errorf("[go-helper] JSON unmarshal failed: %w", err)
	}

	// 如果配置结构体实现了 Validator 接口，执行验证
	if v, ok := any(config.Data).(structs.Validator); ok {
		if e := v.Validate(); e != nil {
			err = fmt.Errorf("[go-helper] Validate config failed, config:%w", e)
			return
		}
	}

	// 如果配置结构体实现了 Formatter 接口，执行格式化
	if v, ok := any(config.Data).(structs.Formatter); ok {
		v.Format()
	}

	// 如果配置结构体实现了 Callbacker 接口，执行回调
	if v, ok := any(config.Data).(structs.Callbacker); ok {
		if e := v.Callback(); e != nil {
			err = fmt.Errorf("[go-helper] Callback config failed, config:%w", e)
			return
		}
	}

	// 监听配置变更
	err = client.ListenConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
		OnChange: func(namespace, group, dataId, data string) {
			// 配置变更时，解析新的配置
			conf := new(T)
			err := json.Unmarshal([]byte(data), conf)
			if err != nil {
				logger.Warn(
					context.TODO(),
					"[go-helper] ConfigOnChange Unmarshal config failed",
					field.String("err", err.Error()),
				)
				return
			}

			// 如果配置结构体实现了 Validator 接口，执行验证
			if v, ok := any(config.Data).(structs.Validator); ok {
				if e := v.Validate(); e != nil {
					logger.Warn(
						context.TODO(),
						"[go-helper] Validate config failed",
						field.String("err", e.Error()),
					)
					return
				}
			}

			// 如果配置结构体实现了 Formatter 接口，执行格式化
			if v, ok := any(config.Data).(structs.Formatter); ok {
				v.Format()
			}

			if v, ok := any(config.Data).(structs.Callbacker); ok {
				if e := v.Callback(); e != nil {
					logger.Warn(
						context.TODO(),
						"[go-helper] Callback config failed",
						field.String("err", e.Error()),
					)
					return
				}
			}

			// 更新配置并记录日志
			config.Set(conf)
			logger.Info(
				context.TODO(),
				"[go-helper] nacos config changed",
				field.String("namespace", namespace),
				field.String("group", group),
				field.String("dataId", dataId),
				field.String("data", data),
			)
		},
	})

	// 如果监听配置时发生错误，记录警告日志
	if err != nil {
		logger.Warn(
			context.TODO(),
			"[go-helper] Get config from nacos failed",
			field.String("dataId", dataId),
			field.String("group", group),
			field.String("err", err.Error()),
		)
		return
	}

	return config, nil
}
