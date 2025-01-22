package define

import "sync"

func NewConfig[T any]() *Config[T] {
	return &Config[T]{
		lock: sync.RWMutex{},
		Data: new(T),
	}
}

type Config[T any] struct {
	lock sync.RWMutex
	Data *T
}

func (config *Config[T]) Get() T {
	config.lock.RLock()
	defer config.lock.RUnlock()
	return *config.Data
}

func (config *Config[T]) Set(data *T) {
	config.lock.Lock()
	defer config.lock.Unlock()
	*config.Data = *data
}
