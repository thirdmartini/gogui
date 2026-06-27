package themes

import (
	"fmt"
	"image"

	"github.com/thirdmartini/gogui/pkg/ux/canvas"
	"github.com/thirdmartini/gogui/pkg/ux/canvas/color"
	"github.com/thirdmartini/gogui/pkg/ux/canvas/fonts"
)

var Default *Theme

func init() {
	Default = &Theme{
		palette: canvas.NewGGCanvas(image.NewRGBA(image.Rect(0, 0, 1, 1))).ColorPalette(),
		//path:       path.Join("assets", "default"),

		searchPaths: []string{"dark", "default"},
		userColors:  NewColorsGroup(),
		fonts:       NewFontCache(),
	}
}

func GetColor(name string) color.Color {
	return Default.GetColor(name)
}

func NewColor(name, hex string) color.Color {
	return Default.NewColor(name, hex)
}

func LoadImage(name string) image.Image {
	return Default.LoadImage(name)
}

func GetIcon(name string) Icon {
	return Default.GetIcon(name)
}

func LoadFont(name string, path string, points float64) (*fonts.Font, error) {
	err := Default.fonts.LoadFont(name, path, points)
	if err != nil {
		return nil, err
	}
	return Default.fonts.Font(name), nil
}

func Font(name string) *fonts.Font {
	if font := Default.fonts.Font(name); font != nil {
		return font
	}
	panic(fmt.Sprintf("trying to load font name(%s) that was not preloaded", name))
}

func SetTheme(theme *Theme) {
	Default = theme
}
