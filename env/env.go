package env

import (
	"os"
	"sync"
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
