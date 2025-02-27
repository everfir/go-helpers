package nacos

import (
	"encoding/json"
	"fmt"

	"github.com/everfir/go-helpers/consts"
	"github.com/everfir/go-helpers/define/config"
	"github.com/everfir/go-helpers/env"
	internal_nacos "github.com/everfir/go-helpers/internal/helper/nacos"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

// NewClient 创建并初始化 Nacos 配置客户端
//
// 该函数执行以下操作：
// 1. 配置客户端参数（超时、日志、缓存等）
// 2. 配置服务器连接信息
// 3. 创建并返回配置客户端实例
//
// 参数：
//   - ip: Nacos 服务器 IP 地址
//   - namespace: 命名空间 ID，用于隔离不同环境的配置
//   - username: Nacos 认证用户名
//   - password: Nacos 认证密码
//
// 返回值：
//   - nacosClient: 初始化成功的 Nacos 配置客户端实例
//   - err: 初始化过程中发生的错误，包括：
//   - 参数校验失败
//   - 服务器连接失败
//   - 认证失败
//
// 示例：
//
//	client, err := NewClient("127.0.0.1", "dev", "nacos", "nacos")
//	if err != nil {
//	    log.Fatal("Failed to create Nacos client:", err)
//	}
func NewClient(ip, namespace, username, password string) (nacosClient config_client.IConfigClient, err error) {
	cc := constant.ClientConfig{
		NamespaceId:         namespace,
		TimeoutMs:           60000,
		NotLoadCacheAtStart: true,
		LogDir:              "",
		CacheDir:            "",
		LogLevel:            "error",
		Username:            username,
		Password:            password,
	}

	sc := []constant.ServerConfig{
		{
			IpAddr: ip,
			Port:   8848,
		},
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"clientConfig":  cc,
		"serverConfigs": sc,
	})
	if err != nil {
		err = fmt.Errorf("[go-helper] nacos.NewClient failed: %w", err)
		return nil, err
	}

	return configClient, nil
}

// GetEverfirNacosClient 获取 Everfir 预配置的 Nacos 客户端实例，用于操作全局业务配置
//
// 该函数返回一个已经初始化好的 Nacos 配置客户端，适用于以下场景：
// 1. 项目中只需要一个全局的 Nacos 客户端
// 2. 客户端配置已经通过环境变量或其他方式预配置
// 3. 需要简化客户端获取逻辑
//
// 返回值：
//   - config_client.IConfigClient: 预配置的 Nacos 客户端实例
//
// 注意：
//  1. 该客户端是单例模式，多次调用返回同一个实例
//  2. 客户端配置通常从环境变量或配置文件加载
//  3. 确保在使用前已经正确配置 Nacos 服务器信息
//
// 示例：
//
//	client := GetEverfirNacosClient()
//	config, err := GetConfigFromNacosAndConfigOnChange[AppConfig](client, "app-config")
//	if err != nil {
//	    log.Fatal(err)
//	}
func GetEverfirNacosClient() config_client.IConfigClient {
	return internal_nacos.GetNacosClient()
}

// GetConfigAndListen 从 Nacos 获取配置并监听配置变更
//
// 该函数执行以下操作：
// 1. 从 Nacos 服务器获取指定 dataId 的配置
// 2. 将配置内容反序列化为指定类型
// 3. 注册配置变更监听器，当配置发生变化时自动更新
//
// 参数：
//   - client: Nacos 配置客户端实例，必须已经初始化
//   - dataId: 配置的唯一标识符
//
// 返回值：
//   - config: 包含配置数据的结构体指针
//   - err: 操作过程中发生的错误，包括：
//   - 配置获取失败
//   - 配置解析失败
//   - 监听器注册失败
//
// 示例：
//
//	type AppConfig struct {
//	    Port int `json:"port"`
//	}
//
//	client := // 初始化 Nacos 客户端
//	cfg, err := GetConfigAndListen[AppConfig](client, "app-config")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("Server port:", cfg.Data.Port)
func GetConfigAndListen[T any](
	client config_client.IConfigClient,
	dataId string,
	gray ...bool,
) (config *config.NacosConfig[T], err error) {
	if len(gray) > 0 && gray[0] {
		return internal_nacos.GetConfigAndListenWithGray[T](client, dataId)
	}
	return internal_nacos.GetConfigAndListen[T](client, dataId)
}

// Publish 向 Nacos 发布配置
//
// 参数：
//   - client: Nacos 配置客户端
//   - dataId: 配置 ID
//   - data: 要发布的配置数据
//   - group: 配置分组，可选参数，默认使用环境变量作为分组
//
// 返回值：
//   - err: 发布过程中发生的错误
//
// 示例：
//
//	err := Publish(client, "my-app-config", EncodingJSON, cfg)
//	if err != nil {
//	    log.Fatal(err)
//	}
func Publish[T any](client config_client.IConfigClient, dataId string, data T, group ...consts.TrafficGroup) (err error) {
	var content string
	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	content = string(b)

	var g string
	if len(group) > 0 {
		g = fmt.Sprintf("%s_%s", env.Env(), group[0].Group())
	} else {
		g = env.Env()
	}

	// 发布配置到 Nacos
	success, err := client.PublishConfig(vo.ConfigParam{
		DataId:  dataId,
		Group:   g,
		Content: content, // 配置内容
		Type:    vo.JSON, // 配置类型
	})
	if err != nil {
		return fmt.Errorf("failed to publish config: %w", err)
	}
	if !success {
		return fmt.Errorf("publish config failed without error")
	}
	return nil
}
