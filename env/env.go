package env

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/everfir/go-helpers/internal/structs"
)

// Env 获取环境变量并设置默认值
//
// 该函数执行以下操作：
// 1. 从环境变量中获取指定键的值
// 2. 如果环境变量为空，默认为测试环境
// 3. 使用 sync.OnceValue 确保只初始化一次
//
// 返回值：
//   - string: 环境变量的值或默认值
//
// 注意：
//  1. 该函数是线程安全的
//  2. 环境变量的键由 EnvKey.String() 提供
//  3. 默认值可以通过修改代码调整
//
// 示例：
//
//	env := Env()
//	if env == EnvTest {
//	    log.Println("running in test environment group")
//	}
var Env func() string = sync.OnceValue[string](func() string {
	env := os.Getenv(EnvKey.String())
	if env == "" {
		env = EnvTest
	}
	return env
})

// Test 判断当前环境是否为测试环境（Test）。
// 该函数通过 Env() 获取当前环境变量，并与 EnvTest 进行比较。
//
// 返回值:
//   - bool: 如果当前环境为测试环境，则返回 true，否则返回 false。
//
// 使用示例:
//
//	if Test() {
//	    fmt.Println("当前是测试环境")
//	}
func Test() bool {
	return Env() == string(EnvTest)
}

// Prod 判断当前环境是否为生产环境（Prod）。
// 该函数通过 Env() 获取当前环境变量，并与 EnvProd 进行比较。
//
// 返回值:
//   - bool: 如果当前环境为生产环境，则返回 true，否则返回 false。
//
// 使用示例:
//
//	if Prod() {
//	    fmt.Println("当前是生产环境")
//	}
func Prod() bool {
	return Env() == string(EnvProd)
}

// Idc 获取当前服务的 IDC（机房）信息。
// 该变量使用 sync.OnceValue 确保 IDC 只计算一次，后续调用将返回缓存的值。
// IDC 信息通过环境变量 "IDC" 获取，若未设置，则默认返回 "IDC_BJ"。
//
// 返回值:
//   - string: 当前服务所在的 IDC 机房名称。
//
// 使用示例:
//
//	fmt.Println(Idc()) // 可能输出: "IDC_BJ" 或环境变量中设置的值
var Idc func() string = sync.OnceValue[string](func() string {
	idc := os.Getenv(IdcKey.String())
	if idc == "" {
		return IDC_BJ
	}
	return idc
})

// CN 判断当前 IDC（Internet Data Center）是否位于北京
// 返回值：
//   - true：当前 IDC 是北京 IDC
//   - false：当前 IDC 不是北京 IDC
//
// 示例：
//
//	if CN() {
//	    fmt.Println("当前 IDC 在北京")
//	}
func CN() bool {
	return Idc() == IDC_BJ
}

// RF 判断当前 IDC（Internet Data Center）是否位于 RF 机房
// 返回值：
//   - true：当前 IDC 是 RF 机房
//   - false：当前 IDC 不是 RF 机房
//
// 示例：
//
//	if RF() {
//	    fmt.Println("当前 IDC 在 RF 机房")
//	}
func RF() bool {
	return Idc() == IDC_RF
}

// Business: 从上下文中获取当前业务
// @param ctx: 上下文
// 前置依赖： middleware.BusinessMiddleware
func Business(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	iface := ctx.Value(BusinessKey)
	business, ok := iface.(string)
	if !ok {
		return ""
	}

	return business
}

// AccountInfo 从给定的上下文（context.Context）中提取账户信息。
// 如果上下文为空或上下文中没有找到有效的账户信息，则返回一个默认的 AccountInfo 结构体。
//
// 参数:
//   - ctx: context.Context，包含账户信息的上下文
//
// 返回值:
//   - info: structs.AccountInfo，从上下文中提取的账户信息
//
// 使用示例:
//
//	ctx := context.WithValue(context.Background(), AccountInfoKey, &structs.AccountInfo{
//	    ID:   "12345",
//	    Name: "Alice",
//	})
//	account := AccountInfo(ctx)
//	if account.Validate() {
//		// 从上下文获取用户信息失败，可能没有调用middleware
//		return
//	}
//	fmt.Println(account.ID, account.Name) // 输出: 12345 Alice
func AccountInfo(ctx context.Context) (info structs.AccountInfo) {
	if ctx == nil {
		return info
	}

	iface := ctx.Value(AccountInfoKey)
	val, ok := iface.(*structs.AccountInfo)
	if !ok {
		return info
	}

	info = *val
	return
}

// Platform 从上下文中获取系统平台信息，返回对应的 TDevicePlatform 类型
func Platform(ctx context.Context) TDevicePlatform {
	if ctx == nil {
		return DP_Unknow
	}

	iface := ctx.Value(PlatformKey)
	platform, ok := iface.(string)
	if !ok {
		return DP_Unknow
	}

	return TDevicePlatform(strings.ToLower(platform))
}

// IOS 判断平台是否为 IOS
func IOS(ctx context.Context) bool {
	return Platform(ctx) == DP_IOS
}

// Android 判断平台是否为 Android
func Android(ctx context.Context) bool {
	return Platform(ctx) == DP_Android
}

// Mac 判断平台是否为 MacOS
func Mac(ctx context.Context) bool {
	return Platform(ctx) == DP_MacOS
}

// Windows 判断平台是否为 Windows
func Windows(ctx context.Context) bool {
	return Platform(ctx) == DP_Windows
}

// Linux 判断平台是否为 Linux
func Linux(ctx context.Context) bool {
	return Platform(ctx) == DP_Linux
}

// Ipad 判断平台是否为 iPadOS
func Ipad(ctx context.Context) bool {
	return Platform(ctx) == DP_IpadOS
}

// Device 从上下文中获取设备信息，返回对应的 Device 类型
func Device(ctx context.Context) TDevice {
	if ctx == nil {
		return Dev_Unknow
	}

	iface := ctx.Value(DeviceKey)
	device, ok := iface.(string)
	if !ok {
		return Dev_Unknow
	}

	return TDevice(strings.ToLower(device))
}

// Phone 判断设备是否为手机
func Phone(ctx context.Context) bool {
	return Device(ctx) == Dev_Phone
}

// PC 判断设备是否为个人电脑
func PC(ctx context.Context) bool {
	return Device(ctx) == Dev_PC
}

// IPad 判断设备是否为 iPad
func IPad(ctx context.Context) bool {
	return Device(ctx) == Dev_IPad
}

// Version 从上下文中获取客户端版本信息
func Version(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	iface := ctx.Value(VersionKey)
	version, ok := iface.(string)
	if !ok {
		return ""
	}

	return version
}

// AppType 根据给定的 context 获取应用类型（app、miniapp 或 web）。
// 它从 context 中提取存储的 AppTypeKey 值，并尝试将其转换为 TAppType。
func AppType(ctx context.Context) TAppType {
	if ctx == nil {
		return ""
	}

	iface := ctx.Value(AppTypeKey)
	appType, ok := iface.(string)
	if !ok {
		return ""
	}

	return TAppType(strings.ToLower(appType))
}

// App 判断应用是否为 app 类型
func App(ctx context.Context) bool {
	return AppType(ctx) == AppType_App
}

// MiniApp 判断应用是否为 miniapp 类型
func MiniApp(ctx context.Context) bool {
	return AppType(ctx) == AppType_MiniApp
}

// Web 判断应用是否为 web 类型
func Web(ctx context.Context) bool {
	return AppType(ctx) == AppType_Web
}

// TestSetAccountInfo:justForTest
func TestSetAccountInfo(ctx context.Context, id uint64) context.Context {
	if !Test() {
		return ctx
	}
	return context.WithValue(ctx, AccountInfoKey, &structs.AccountInfo{AccountId: id, TemplateIDs: []string{fmt.Sprintf("%d", id)}})
}
