package themes

import (
	"fmt"
	"image"
	"path"

	"github.com/thirdmartini/gogui/pkg/ux/canvas/color"
	"github.com/thirdmartini/gogui/pkg/ux/canvas/fonts"
)

var Default *Theme

func init() {
	Default = &Theme{
		Path:       path.Join("assets", "dark"),
		Colors:     &SystemColors{},
		UserColors: make(map[string]color.Color),
	}
}

func GetColor(name string) color.Color {
	return Default.GetColor(name)
}

func LoadImage(name string) image.Image {
	return Default.LoadImage(name)
}

func LoadFont(name string, path string, points float64) (*fonts.Font, error) {
	err := Default.Fonts.LoadFont(name, path, points)
	if err != nil {
		return nil, err
	}
	return Default.Fonts.Font(name), nil
}

func Font(name string) *fonts.Font {
	if font := Default.Fonts.Font(name); font != nil {
		return font
	}
	panic(fmt.Sprintf("trying to load font name(%s) that was not preloaded", name))
}

func SetTheme(theme *Theme) {
	Default = theme
}
