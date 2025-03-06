package consts

type IDC uint8

const (
	IDC_BJ string = "beijing"
	IDC_SH string = "shanghai"
	IDC_GZ string = "guangzhou"
	IDC_RF string = "rf"
)

const (
	EnvTest string = "test"
	EnvProd        = "production"
)

type ContextKey string

func (key ContextKey) String() string {
	return string(key)
}

const (
	// EnvKey: 环境变量中标识服务所处环境
	EnvKey ContextKey = "ENV"
	// IdcKey: 环境变量中标识服务所处idc
	IdcKey ContextKey = "IDC"

	// DeviceKey: 请求头&上下文中携带设备信息
	DeviceKey ContextKey = "x-everfir-device"
	// VersionKey: 请求头&上下文中携带客户端版本信息
	VersionKey ContextKey = "x-everfir-version"
	// AppTypeKey: 请求头&上下文中携带用户App类型信息
	AppTypeKey ContextKey = "x-everfir-app-type"
	// PlatformKey: 请求头&上下文中携带平台信息
	PlatformKey ContextKey = "x-everfir-platform"
	// BusinessKey: 请求头中携带业务信息
	BusinessKey ContextKey = "x-everfir-business"
	// AccountInfoKey: 用户信息，请求头&上下文中携带用户信息
	AccountInfoKey ContextKey = "x-everfir-account-info"
	// ExperimentGroupKey: 请求头中携带分组信息, 用于AB分组
	ExperimentGroupKey ContextKey = "x-everfir-experiment-group"
)
