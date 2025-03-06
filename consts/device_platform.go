package consts

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
