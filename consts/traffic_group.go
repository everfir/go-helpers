package consts

func NewTrafficGroup(group uint8) TrafficGroup {
	var g = TrafficGroup(group)
	if g < TrafficGroup_A || g >= TrafficGroup_Unknow {
		return TrafficGroup_Unknow
	}

	return g
}

func NewTrafficGroupFromString(group string) TrafficGroup {
	if len(group) != 1 {
		return TrafficGroup_A
	}

	var g = TrafficGroup(group[0] - 'a')
	if g < TrafficGroup_A || g >= TrafficGroup_Unknow {
		return TrafficGroup_A
	}

	return g
}

type TrafficGroup uint8

const (
	TrafficGroup_A      TrafficGroup = iota // 0
	TrafficGroup_B                          // 1
	TrafficGroup_C                          // 2
	TrafficGroup_D                          // 3
	TrafficGroup_E                          // 4
	TrafficGroup_F                          // 5
	TrafficGroup_G                          // 6
	TrafficGroup_H                          // 7
	TrafficGroup_I                          // 8
	TrafficGroup_J                          // 9
	TrafficGroup_K                          // 10
	TrafficGroup_L                          // 11
	TrafficGroup_M                          // 12
	TrafficGroup_N                          // 13
	TrafficGroup_O                          // 14
	TrafficGroup_P                          // 15
	TrafficGroup_Q                          // 16
	TrafficGroup_R                          // 17
	TrafficGroup_S                          // 18
	TrafficGroup_T                          // 19
	TrafficGroup_U                          // 20
	TrafficGroup_V                          // 21
	TrafficGroup_W                          // 22
	TrafficGroup_X                          // 23
	TrafficGroup_Y                          // 24
	TrafficGroup_Z                          // 25
	TrafficGroup_Unknow                     // 26
)

func (g TrafficGroup) Group() string {
	if g <= TrafficGroup_A || g >= TrafficGroup_Unknow {
		return "a"
	}

	return string('a' + g)
}
