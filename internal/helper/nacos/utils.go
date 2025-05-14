package nacos

import (
	"os"
	"sync"

	"github.com/everfir/go-helpers/consts"
	"github.com/everfir/go-helpers/env"
)

var ipMapping func() map[string]string = sync.OnceValue(func() map[string]string {
	switch env.Idc() {
	case consts.IDC_BJ:
		return map[string]string{
			// env.EnvTest: "192.168.0.8",
			consts.EnvTest: "101.126.144.112",
			consts.EnvProd: "192.168.0.49",
		}
	case consts.IDC_RF:
		return map[string]string{
			consts.EnvTest: "192.168.0.8",
			consts.EnvProd: "192.168.0.49",
		}
	default:
		return map[string]string{}
	}
})

var namespaceMapping func() map[string]string = sync.OnceValue(func() map[string]string {
	switch env.Idc() {
	case consts.IDC_BJ:
		return map[string]string{
			consts.EnvTest: "56240543-0336-4fe4-815d-d2437c2bb11e",
			consts.EnvProd: "a5299e86-dbe0-409a-bdd9-7a7e6ef346ba",
		}
	case consts.IDC_RF:
		return map[string]string{
			consts.EnvTest: "93e786d8-09d5-4106-a99c-2eee98f707b6",
			consts.EnvProd: "d80b88e0-d583-44e3-92f8-4b93164a92c7",
		}
	default:
		return map[string]string{}
	}
})

var authMapping func() map[string][]string = sync.OnceValue(func() map[string][]string {
	switch env.Idc() {
	case consts.IDC_BJ:
		return map[string][]string{
			consts.EnvTest: []string{"nacos", "EverFir@Nacos20240717.."},
			consts.EnvProd: []string{"nacos", `BNg%d59=`},
		}
	case consts.IDC_RF:
		return map[string][]string{
			consts.EnvTest: []string{"nacos", "5oadA-c)"},
			consts.EnvProd: []string{"nacos", "*60gdE8q"},
		}
	default:
		return map[string][]string{}
	}
})

func NacosIp() string {
	ret := os.Getenv("EVERFIR_NACOS_IP")
	if ret == "" {
		ret = ipMapping()[env.Env()]
	}
	return ret
}

func Namespace() string {
	ret := os.Getenv("EVERFIR_NACOS_NAMESPACE")
	if ret == "" {
		ret = namespaceMapping()[env.Env()]
	}
	return ret
}

func AuthInfo() (username, passward string) {
	info := authMapping()[env.Env()]
	if len(info) == 0 {
		return
	}
	return info[0], info[1]
}
