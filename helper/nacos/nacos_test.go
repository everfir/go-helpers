package nacos_test

import (
	"testing"
	"time"

	"github.com/everfir/go-helpers/helper/nacos"
	"github.com/everfir/go-helpers/internal/structs"
	"github.com/stretchr/testify/assert"
)

func TestNacosConfig(t *testing.T) {
	// 初始化Nacos客户端
	client, err := nacos.NewClient("101.126.144.112", "56240543-0336-4fe4-815d-d2437c2bb11e", "nacos", "EverFir@Nacos20240717..")
	if err != nil {
		t.Fatalf("Failed to create Nacos client: %v", err)
	}

	// 测试用例1：获取配置
	t.Run("GetConfig", func(t *testing.T) {
		var config *structs.Config[map[string]bool]
		config, err := nacos.GetConfigFromNacosAndConfigOnChange[map[string]bool](client, "shutdown.json")
		if err != nil {
			t.Fatalf("Failed to get config: %v", err)
		}

		t.Logf("Initial config: %+v", config.Data)
		assert.NotNil(t, config.Data, "Config data should not be nil")
	})

	// 测试用例2：配置更新监听
	t.Run("ConfigUpdate", func(t *testing.T) {
		var config *structs.Config[map[string]bool]
		config, err := nacos.GetConfigFromNacosAndConfigOnChange[map[string]bool](client, "shutdown.json")
		if err != nil {
			t.Fatalf("Failed to get config: %v", err)
		}

		// 模拟配置更新
		newConfig := map[string]bool{
			"shutdown.enabled": true,
			"shutdown.timeout": false,
		}
		err = nacos.Publish(client, "shutdown.json", newConfig)
		if err != nil {
			t.Fatalf("Failed to publish config: %v", err)
		}

		// 等待配置更新
		time.Sleep(2 * time.Second)

		t.Logf("Updated config: %+v", config.Data)
		assert.True(t, config.Get()["shutdown.enabled"], "Shutdown enabled should be true")
		assert.False(t, config.Get()["auto_restart"], "Auto restart should be false")
	})

	// 测试用例3：错误场景
	t.Run("ErrorCases", func(t *testing.T) {
		// 测试不存在的配置
		_, err := nacos.GetConfigFromNacosAndConfigOnChange[map[string]bool](client, "non_existent.json")
		assert.Error(t, err, "Should return error for non-existent config")

		// 测试错误的dataId格式
		_, err = nacos.GetConfigFromNacosAndConfigOnChange[map[string]bool](client, "")
		assert.Error(t, err, "Should return error for empty dataId")
	})
}
