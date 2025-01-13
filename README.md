# Go Helpers

这是一个基于Gin框架的Go语言工具库，主要用于处理服务优雅关闭和配置管理。项目集成了Nacos配置中心，提供了灵活的中间件机制，特别适用于微服务架构下的服务管理。

## 功能特性

- **优雅关闭**：通过中间件实现服务的优雅关闭
- **配置管理**：集成Nacos配置中心，支持动态配置更新
- **请求过滤**：基于业务标识的请求过滤机制
- **状态监控**：提供简单的服务状态监控接口

## 快速开始

### 安装依赖
```bash
go get github.com/everfir/go-helpers
```

### 运行示例
```bash
go run internal/example/shutdown.go
```

## 项目结构
.
├── env # 集群环境识别工具
├── internal
│ ├── example
│ │ └── shutdown.go # 示例代码
│ ├── helper
│ │ └── nacos
│ │ ├── nacos.go # Nacos配置管理
│ │ └── nacos_test.go # Nacos测试
├── middleware
│ └── shutdown_middleware.go # 停服
├── go.mod
├── go.sum
└── README.md

## 配置说明
项目使用Nacos作为配置中心，配置文件格式如下：
### 停服
```json
// shutdown.json
{
    "business1": true,
    "business2": false
}
```
