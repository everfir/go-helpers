package account

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/everfir/go-helpers/consts"
	"github.com/everfir/go-helpers/define"
	"github.com/everfir/go-helpers/define/config"
	"github.com/everfir/go-helpers/internal/helper/nacos"

	"github.com/everfir/go-helpers/env"
	util_http "github.com/everfir/go-helpers/internal/util/http"
)

var getAccountConfig func() *config.NacosConfig[define.AccountConfig] = sync.OnceValue(func() *config.NacosConfig[define.AccountConfig] {
	config, err := nacos.GetConfigAndListen[define.AccountConfig](nacos.GetNacosClient(), "account_config.json")
	if err != nil {
		panic(err.Error())
	}

	return config
})

//const (
//	checkTokenUrl     string = "http://user-account:8080/account/check_token"
//	testCheckTokenUrl string = "http://101.126.81.38:10003/account/check_token"
//)

func getTokenUrl() string {
	cfg, _ := getAccountConfig().Get(consts.TrafficGroup_B)

	if env.Prod() {
		return cfg.UrlEnv[consts.EnvProd]
	}

	cfg, _ = getAccountConfig().Get()
	return cfg.UrlEnv[consts.EnvTest]
}

type CheckTokenReq struct {
	Token string `json:"token,omitempty"`
}

type CheckTokenResp struct {
	AccountInfo *define.AccountInfo `json:"account_info,omitempty"`
	Valid       bool                `json:"valid,omitempty"`
	ErrCode     uint32              `json:"err_code,omitempty"`
	ErrMsg      string              `json:"err_msg,omitempty"`
}

var (
	client = http.Client{
		Timeout:   5 * time.Second,
		Transport: util_http.NewTraceTripper(http.DefaultTransport),
	}
)

func CheckToken(ctx context.Context, token string) (*CheckTokenResp, error) {
	req := CheckTokenReq{
		Token: token,
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, getTokenUrl(), bytes.NewReader(reqBytes))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	all, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	resp := &CheckTokenResp{}
	err = json.Unmarshal(all, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
