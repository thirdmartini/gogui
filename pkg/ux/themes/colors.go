package themes

import (
	"fmt"
	"image"

	"github.com/thirdmartini/gogui/pkg/ux/canvas"
	"github.com/thirdmartini/gogui/pkg/ux/canvas/color"
)

/*
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
)*/

type SystemColors struct {
	None       color.Color
	Foreground color.Color
	Background color.Color
	Border     color.Color

	TextPrimary color.Color
	TextMuted   color.Color

	MenuBackground   color.Color
	WindowBackground color.Color
}

var Colors = SystemColors{}

//var UserColors UserColorGroup

func initSystemColors() {
	// this is hackery for now
	palette := canvas.NewGGCanvas(image.NewRGBA(image.Rect(0, 0, 1, 1))).ColorPalette()
	Colors = SystemColors{
		Foreground: palette.NewRGB8(255, 255, 255),
		Background: palette.NewRGB8(0, 0, 0),
		Border:     palette.NewRGB8(222, 222, 222),

		TextPrimary: palette.NewRGB8(255, 255, 255),
		TextMuted:   palette.NewRGB8(128, 128, 128),

		MenuBackground:   palette.NewRGB8(0, 0, 0),
		WindowBackground: palette.NewRGB8(0, 0, 0),
	}
	Default.Colors = &Colors

}

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

func toColor(p color.Palette, hexColor string) color.Color {
	var r, g, b uint8

	if cnt, err := fmt.Sscanf(hexColor, "#%02x%02x%02x", &r, &g, &b); err == nil && cnt == 3 {
		return p.NewRGB8(r, g, b)
	}

	return p.NewRGB8(0, 0, 0)
}

type UserColorGroup map[string]color.Color

func (ug UserColorGroup) GetColor(name string) color.Color {
	col := ug[name]
	return col
}
