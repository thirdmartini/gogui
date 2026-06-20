package composer

import (
	"github.com/thirdmartini/gogui/pkg/ux"
)

var AlignValueStrings = map[string]int{
	"":        ux.AlignDefault,
	"default": ux.AlignDefault,
	"left":    ux.AlignLeft,
	"right":   ux.AlignRight,
	"center":  ux.AlignCenter,
}

func alignValue(s string) uint8 {
	v, _ := AlignValueStrings[s]
	return uint8(v)
}
