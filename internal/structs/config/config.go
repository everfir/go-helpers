package config

import "sync"

func NewConfig[T any]() *Config[T] {
	return &Config[T]{
		lock:      sync.RWMutex{},
		Data:      new(T),
		listeners: map[string]IListener[T]{},
	}
}

type IListener[T any] interface {
	OnChange(data T)
}

type Config[T any] struct {
	lock      sync.RWMutex
	Data      *T
	listeners map[string]IListener[T]
}

func (config *Config[T]) Get() T {
	config.lock.RLock()
	defer config.lock.RUnlock()

	return *config.Data
}

func (config *Config[T]) Set(data *T) {
	config.lock.Lock()
	config.Data = data
	config.lock.Unlock()

	for _, listener := range config.listeners {
		listener.OnChange(*data)
	}
}

func (config *Config[T]) RegisterListener(name string, listener IListener[T]) {
	config.lock.Lock()
	defer config.lock.Unlock()
	config.listeners[name] = listener
}

func (config *Config[T]) UnregisterListener(name string) {
	config.lock.Lock()
	defer config.lock.Unlock()

	delete(config.listeners, name)
}
