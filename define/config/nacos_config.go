package config

import (
	"fmt"
	"sync"

	"github.com/everfir/go-helpers/consts"
	"github.com/everfir/go-helpers/env"
	internal_config "github.com/everfir/go-helpers/internal/structs/config"
)

func NewNacosConfig[V any](data map[string]*internal_config.Config[V]) *NacosConfig[V] {
	return &NacosConfig[V]{
		lock: sync.RWMutex{},
		data: data,
	}
}

type NacosConfig[V any] struct {
	lock sync.RWMutex
	data map[string]*internal_config.Config[V]
}

// Get 获取配置数据，支持按流量分组获取不同的配置。
//
// 该方法会根据环境和流量分组来获取对应的配置数据：
// 1. 如果指定了流量分组，会尝试获取 "{env}_{group}" 格式的配置。
// 2. 如果未找到对应分组的配置，会降级使用 A 组（稳定组）的配置。
// 3. 如果未指定流量分组，则直接使用当前环境的默认配置。
//
// 参数：
//   - keys: 可选的流量分组参数，如果提供则按指定分组获取配置。
//
// 返回值：
//   - V: 泛型类型的配置数据。
//
// 示例：
//
//	config.Get()                    // 获取当前环境的默认配置
//	config.Get(TrafficGroup_B)      // 获取 B 组的配置，如果不存在则返回 A 组配置
func (config *NacosConfig[V]) Get(keys ...consts.TrafficGroup) (V, bool) {
	// 加读锁保护并确保解锁
	config.lock.RLock()
	defer config.lock.RUnlock()

	// 获取当前环境作为默认 key
	k := env.Env()

	// 如果指定了流量分组，将环境和分组组合成新的 key
	if len(keys) > 0 {
		k = fmt.Sprintf("%s_%s", k, keys[0].Group())
	}

	var exist bool = true
	// 如果找不到对应分组的配置，默认使用 A 组配置
	if _, exist = config.data[k]; !exist {
		k = fmt.Sprintf("%s_%s", env.Env(), consts.TrafficGroup_A.Group())
	}

	// 返回对应的配置数据
	return config.data[k].Get(), exist
}

// RegisterListener 为指定分组的内部配置注册监听器。
//
// 用途：
// - 当某个分组（例如环境分组或灰度分组）的配置发生变更时，收到回调通知。
//
// 参数：
// - name: 分组键，同时也作为内部 Config 的监听器名称键。
//   - 非灰度场景：通常为当前环境，如 "test"、"prod"。
//   - 灰度场景：通常为形如 "{env}_{group}" 的完整分组键，如 "test_A"、"test_B"。
//     必须确保该分组键已在 `NacosConfig.data` 中存在，否则将导致空指针异常。
//
// - listener: 实现了 `internal_config.IListener[V]` 的监听器实例，会在目标分组配置更新后接收最新的数据副本。
//
// 并发：
// - 本方法内部加写锁，线程安全；监听器回调由内部 `Config.Set` 触发。
func (config *NacosConfig[V]) RegisterListener(name string, listener internal_config.IListener[V]) {
	config.lock.Lock()
	defer config.lock.Unlock()

	config.data[name].RegisterListener(name, listener)
}

// UnregisterListener 从指定分组的内部配置中注销监听器。
//
// 参数：
//   - name: 分组键，同时也是注册时使用的监听器名称键。
//     必须确保该分组键已在 `NacosConfig.data` 中存在，否则将导致空指针异常。
//
// 行为：
// - 如果对应名称的监听器不存在，内部删除操作为幂等，不会报错。
//
// 并发：
// - 本方法内部加写锁，线程安全。
func (config *NacosConfig[V]) UnregisterListener(name string) {
	config.lock.Lock()
	defer config.lock.Unlock()

	config.data[name].UnregisterListener(name)
}
