package nacos_test

import (
	"testing"

	"github.com/everfir/go-helpers/internal/helper/nacos"
)

func TestNacos(t *testing.T) {

	var config *nacos.Config[map[string]string]
	config, err := nacos.GetConfigFromNacosAndConfigOnChange[map[string]string]("shutdown.json")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("getConfig From Nacos: %v", config.Data)
}
