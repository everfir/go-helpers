package nacos

import (
	"sync"

	"github.com/everfir/go-helpers/env"
)

var ipMapping func() map[string]string = sync.OnceValue(func() map[string]string {
	switch env.Idc() {
	case env.IDC_BJ:
		return map[string]string{
			// env.EnvTest: "192.168.0.8",
			env.EnvTest: "101.126.144.112",
			env.EnvProd: "192.168.0.49",
		}
	case env.IDC_RF:
		return map[string]string{
			env.EnvTest: "192.168.0.8",
			env.EnvProd: "192.168.0.49",
		}
	default:
		return map[string]string{}
	}
})

var namespaceMapping func() map[string]string = sync.OnceValue(func() map[string]string {
	switch env.Idc() {
	case env.IDC_BJ:
		return map[string]string{
			env.EnvTest: "56240543-0336-4fe4-815d-d2437c2bb11e",
			env.EnvProd: "a5299e86-dbe0-409a-bdd9-7a7e6ef346ba",
		}
	case env.IDC_RF:
		return map[string]string{
			env.EnvTest: "93e786d8-09d5-4106-a99c-2eee98f707b6",
			env.EnvProd: "d80b88e0-d583-44e3-92f8-4b93164a92c7",
		}
	default:
		return map[string]string{}
	}
})

var authMapping func() map[string][]string = sync.OnceValue(func() map[string][]string {
	switch env.Idc() {
	case env.IDC_BJ:
		return map[string][]string{
			env.EnvTest: []string{"nacos", "EverFir@Nacos20240717.."},
			env.EnvProd: []string{"nacos", `BNg%d59=`},
		}
	case env.IDC_RF:
		return map[string][]string{
			env.EnvTest: []string{"nacos", "5oadA-c)"},
			env.EnvProd: []string{"nacos", "*60gdE8q"},
		}
	default:
		return map[string][]string{}
	}
})

func NacosIp() string {
	return ipMapping()[env.Env()]
}

func Namespace() string {
	return namespaceMapping()[env.Env()]
}

func AuthInfo() (username, passward string) {
	info := authMapping()[env.Env()]
	if len(info) == 0 {
		return
	}
	return info[0], info[1]
}
