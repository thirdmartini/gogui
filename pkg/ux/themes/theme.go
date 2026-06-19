package themes

import (
	"encoding/json"
	"image"
	"os"
	"path"

	"github.com/thirdmartini/gogui/pkg/ux/canvas"
	"github.com/thirdmartini/gogui/pkg/ux/canvas/color"
)

var ThemePath = path.Join("assets", "dark")

type Colors struct {
	Foreground color.Color
	Background color.Color
	Border     color.Color

	TextPrimary color.Color
	TextMuted   color.Color

	MenuBackground   color.Color
	WindowBackground color.Color
}

type Theme struct {
	Colors Colors
}

var Default = Theme{}

func init() {
	palette := canvas.NewGGCanvas(image.NewRGBA(image.Rect(0, 0, 1, 1))).ColorPalette()

	Default.Colors = Colors{
		Foreground: palette.NewRGB8(255, 255, 255),
		Background: palette.NewRGB8(0, 0, 0),
		Border:     palette.NewRGB8(222, 222, 222),

		TextPrimary: palette.NewRGB8(255, 255, 255),
		TextMuted:   palette.NewRGB8(128, 128, 128),

		MenuBackground:   palette.NewRGB8(0, 0, 0),
		WindowBackground: palette.NewRGB8(0, 0, 0),
	}
}

func SetTheme(path string) {
	ThemePath = path
}

func LoadColors(palette color.Palette) error {
	data, err := os.ReadFile(path.Join(ThemePath, "colors.json"))
	if err == nil {
		tc := new(ThemeColors)
		err = json.Unmarshal(data, tc)
		Default.Colors.Background = tc.toColor(palette, "background")
		Default.Colors.Foreground = tc.toColor(palette, "foreground")
		Default.Colors.TextPrimary = tc.toColor(palette, "text:primary")
		Default.Colors.TextMuted = tc.toColor(palette, "text:muted")
		Default.Colors.WindowBackground = tc.toColor(palette, "window:background")
		Default.Colors.Border = tc.toColor(palette, "border")
		Default.Colors.MenuBackground = tc.toColor(palette, "menu:background")
		return nil
	}
	return err
}
