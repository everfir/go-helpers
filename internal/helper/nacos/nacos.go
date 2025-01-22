package nacos

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/everfir/go-helpers/env"
	"github.com/everfir/go-helpers/internal/structs"
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

func GetConfigFromNacosAndConfigOnChange[T any](client config_client.IConfigClient, dataId string) (config *structs.Config[T], err error) {
	cfg, err := client.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  env.Env(),
	})
	if err != nil {
		err = fmt.Errorf("[go-helper] Get config from nacos failed, err: %w", err)
		return
	}

	config = structs.NewConfig[T]()
	err = json.Unmarshal([]byte(cfg), config.Data)
	if err != nil {
		return nil, fmt.Errorf("[go-helper] JSON unmarshal failed: %w", err)
	}

	if v, ok := any(config.Data).(structs.Validator); ok {
		if !v.Validate() {
			err = fmt.Errorf("[go-helper] Validate config failed, config:%+v", config.Data)
			return
		}
	}

	if v, ok := any(config.Data).(structs.Formatter); ok {
		v.Format()
	}

	err = client.ListenConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  env.Env(),
		OnChange: func(namespace, group, dataId, data string) {
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
	if err != nil {
		logger.Warn(
			context.TODO(),
			"[go-helper] Get config from nacos failed",
			field.String("dataId", dataId),
			field.String("group", env.Env()),
			field.String("err", err.Error()),
		)
	}
	return config, nil
}
