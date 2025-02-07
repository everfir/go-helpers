package external_api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/everfir/go-helpers/define"
	"github.com/everfir/go-helpers/env"
	"github.com/everfir/go-helpers/internal/helper/nacos"
	"github.com/everfir/go-helpers/internal/structs"
	"io"
	"net/http"
	"sync"
)

var GetAccountCfg func() *define.Config[structs.AccountServiceConfig] = sync.OnceValue(func() *define.Config[structs.AccountServiceConfig] {
	config, err := nacos.GetConfigFromNacosAndConfigOnChange[structs.AccountServiceConfig](nacos.GetNacosClient(), "account_service_config.json")
	if err != nil {
		panic(err.Error())
	}
	return config
})

func CheckToken(ctx context.Context, token, routerKey string) (*structs.CheckTokenResp, string, error) {
	req := structs.CheckTokenReq{
		Token: token,
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, "", err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, GetAccountCfg().Get().Url, bytes.NewReader(reqBytes))
	if err != nil {
		return nil, "", err
	}

	request.Header.Set("Content-Type", "application/json")

	request.Header.Set(string(env.RouterKey), routerKey)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, "", err
	}

	routerGroup := response.Header.Get(string(env.RouterGroupKey))

	all, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, "", err
	}

	resp := &structs.CheckTokenResp{}
	err = json.Unmarshal(all, resp)
	if err != nil {
		return nil, "", err
	}

	return resp, routerGroup, nil
}
