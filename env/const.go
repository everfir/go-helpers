package env

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

	// RouterKey: 请求头中携带路由信息, 用于AB分组
	RouterKey ContextKey = "x-everifr-router"
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
	// ExperimentGroupKey: 请求头中携带实验分组信息
	ExperimentGroupKey ContextKey = "x-everfir-experiment-group"
)

type TDevicePlatform string

func (p TDevicePlatform) String() string {
	return string(p)
}

const (
	DP_IOS     TDevicePlatform = "ios"
	DP_Linux   TDevicePlatform = "linux"
	DP_MacOS   TDevicePlatform = "macos"
	DP_IpadOS  TDevicePlatform = "ipados"
	DP_Windows TDevicePlatform = "windows"
	DP_Android TDevicePlatform = "android"
	DP_Unknow  TDevicePlatform = ""
)

type TDevice string

func (d TDevice) String() string {
	return string(d)
}

const (
	Dev_Phone  TDevice = "phone"
	Dev_PC     TDevice = "pc"
	Dev_IPad   TDevice = "ipad"
	Dev_Unknow TDevice = ""
)

type TAppType string

func (a TAppType) String() string {
	return string(a)
}

const (
	AppType_App     TAppType = "app"
	AppType_MiniApp TAppType = "miniapp"
	AppType_Web     TAppType = "web"
	AppType_Unknow  TAppType = "" // 可选的，表示未知类型
)
