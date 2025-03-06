package nacos_test

import (
	"context"
	"testing"
	"time"

	"github.com/everfir/go-helpers/consts"
	"github.com/everfir/go-helpers/helper/nacos"
	"github.com/everfir/logger-go"
	"github.com/everfir/logger-go/structs/field"
	"github.com/everfir/logger-go/structs/log_level"
)

// TestGetConfigAndListen 测试从 Nacos 获取配置并监听配置变更
func TestGetConfigAndListen(t *testing.T) {
	logger.Init(logger.WithLevel(log_level.InfoLevel))
	client := nacos.GetEverfirNacosClient()

	cfg, err := nacos.GetConfigAndListen[map[string]bool](client, "shutdown.json", true)
	if err != nil {
		t.Fatalf("Failed to get config from Nacos: %v", err)
	}

	// 前提： 需要先配置好配置

	cfgDetail, exist := cfg.Get()
	logger.Info(context.Background(), "A组配置 ", field.Any("config", cfgDetail), field.Bool("exist", exist))
	// 修改配置, 观察日志
	nacos.Publish(client, "shutdown.json", map[string]bool{
		"a": true,
	})

	cfgDetail, exist = cfg.Get()
	logger.Info(context.Background(), "B组配置 ", field.Any("config", cfgDetail), field.Bool("exist", exist))
	// 修改配置, 观察日志
	nacos.Publish(client, "shutdown.json", map[string]bool{
		"b": true,
	}, consts.TrafficGroup_B)

	// 修改配置, 观察日志
	cfgDetail, exist = cfg.Get(consts.TrafficGroup_Z)
	logger.Info(context.Background(), "Z组配置 ", field.Any("config", cfgDetail), field.Bool("exist", exist))
	nacos.Publish(client, "shutdown.json", map[string]bool{
		"z": true,
	}, consts.TrafficGroup_Z)

	logger.Info(context.Background(), "等待10s后结束")
	for i := 0; i < 3; i++ {
		time.Sleep(1 * time.Second)
	}

	cfgDetail, exist = cfg.Get(consts.TrafficGroup_B)
	logger.Info(context.Background(), "修改之后的A组配置 ", field.Any("config", cfgDetail), field.Bool("exist", exist))
	cfgDetail, exist = cfg.Get(consts.TrafficGroup_B)
	logger.Info(context.Background(), "修改之后的B组配置 ", field.Any("config", cfgDetail), field.Bool("exist", exist))
	cfgDetail, exist = cfg.Get(consts.TrafficGroup_Z)
	logger.Info(context.Background(), "修改之后的Z组配置 ", field.Any("config", cfgDetail), field.Bool("exist", exist))

	logger.Info(context.Background(), "test done")
}
