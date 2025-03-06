package consts

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
