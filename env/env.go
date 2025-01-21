package env

import (
	"context"
	"os"
	"sync"
)

type ContextKey string

func (key ContextKey) String() string {
	return string(key)
}

const (
	// BusinessKey: 请求中携带业务信息的Header
	BusinessKey ContextKey = "x-everfir-business"
	// RouterKey: 请求中携带路由信息的Header, 用于AB分组
	RouterKey ContextKey = "x-everifr-router"
	// RouterGroupKey: 请求/响应中携带路由分组信息的Header
	RouterGroupKey ContextKey = "x-everfir-router-group"
	// ExperimentGroupKey: 请求中携带实验分组信息的Header
	ExperimentGroupKey ContextKey = "x-everfir-experiment-group"
)

var Env func() string = sync.OnceValue[string](func() string {
	env := os.Getenv("NODE_ENV")
	if env == "" {
		env = "DEFAULT_GROUP"
	}
	return env
})

func Test() bool {
	return Env() == string(EnvTest)
}

func Prod() bool {
	return Env() == string(EnvProd)
}

var Idc func() string = sync.OnceValue[string](func() string {
	idc := os.Getenv("IDC")
	if idc == "" {
		return IDC_BJ
	}
	return idc
})

func CN() bool {
	return Idc() == IDC_BJ
}

func RF() bool {
	return Idc() == IDC_RF
}

// Business: 从上下文中获取当前业务
func Business(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	iface := ctx.Value(BusinessKey)
	business, ok := iface.(string)
	if !ok {
		return ""
	}

	return business
}

func ExperimentGroup(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	iface := ctx.Value(ExperimentGroupKey)
	group, ok := iface.(string)
	if !ok {
		return ""
	}

	return group
}
