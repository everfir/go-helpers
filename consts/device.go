package consts

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
