package gray

import (
	"context"
	"sync"

	"github.com/everfir/go-helpers/env"
	"github.com/everfir/go-helpers/internal/helper/nacos"
	"github.com/everfir/go-helpers/internal/structs"
)

var grayConfig func() *structs.Config[structs.Gray] = sync.OnceValue(func() *structs.Config[structs.Gray] {
	config, err := nacos.GetConfigFromNacosAndConfigOnChange[structs.Gray]("gray.json")
	if err != nil {
		panic(err.Error())
	}

	return config
})

func Gray(ctx context.Context, feature, user string) bool {
	business := env.Business(ctx)
	if business == "" {
		return false
	}

	return grayConfig().Get().Enable(business, feature, user)
}
