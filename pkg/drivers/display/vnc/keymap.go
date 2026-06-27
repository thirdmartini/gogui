package vnc

import (
	"github.com/thirdmartini/gogui/pkg/ux"
)

const (
	KeyboardPressUp    = 0xff52
	KeyboardPressDown  = 0xff54
	KeyboardPressLeft  = 0xff51
	KeyboardPressRight = 0xff53
)

var remap = map[uint32]uint{
	KeyboardPressUp:    ux.KeyPressUp,
	KeyboardPressDown:  ux.KeyPressDown,
	KeyboardPressLeft:  ux.KeyPressLeft,
	KeyboardPressRight: ux.KeyPressRight,
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
