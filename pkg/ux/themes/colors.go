package themes

import (
	"fmt"

	"github.com/thirdmartini/gogui/pkg/ux/canvas/color"
)

var (
	ColorBackground color.Color
	ColorBorder     color.Color
	ColorForeground color.Color

	ColorTextPrimary color.Color
	ColorTextMuted   color.Color
	ColorGraphAxis   []color.Color
	ColorGraphTicks  color.Color

	ColorMenuBackground color.Color

	ColorWindowBackground color.Color
)

type ThemeColors map[string]string

func (tc ThemeColors) toColor(p color.Palette, c string) color.Color {
	var r, g, b uint8

	col, ok := tc[c]
	if ok {
		if cnt, err := fmt.Sscanf(col, "#%02x%02x%02x", &r, &g, &b); err == nil && cnt == 3 {
			fmt.Printf("Loaded color [%s] = [%d, %d, %d]\n", c, r, g, b)
			return p.NewRGB8(r, g, b)
		}
	}

	fmt.Printf("Could not load color for [%s]\n", c)

	return p.NewRGB8(0, 0, 0)
}
