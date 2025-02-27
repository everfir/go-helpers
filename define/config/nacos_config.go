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

func (config *NacosConfig[V]) Get(keys ...consts.TrafficGroup) V {
	config.lock.RLock()
	defer config.lock.RUnlock()

	k := env.Env()
	if len(keys) > 0 {
		k = fmt.Sprintf("%s_%s", k, keys[0].Group())
	}

	return config.data[k].Get()
}
