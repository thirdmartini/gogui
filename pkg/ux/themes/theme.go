package themes

import (
	"encoding/json"
	"image"
	"os"
	"path"
	"strings"

	"github.com/thirdmartini/gogui/pkg/ux/canvas/color"
)

type Theme struct {
	Path       string
	Colors     *SystemColors
	UserColors UserColorGroup
	Fonts      *FontCache
}

func (t *Theme) GetColor(name string) color.Color {
	return t.UserColors.GetColor(name)
}

func (t *Theme) LoadImage(name string) image.Image {
	imgPath := path.Join(t.Path, name)
	infile, err := os.Open(imgPath)
	if err != nil {
		return nil
	}
	defer infile.Close()

	// Decode will figure out what type of image is in the file on its own.
	// We just have to be sure all the image packages we want are imported.
	src, _, err := image.Decode(infile)
	if err != nil {
		return nil
	}
	return src
}

func Load(themePath string, palette color.Palette) (*Theme, error) {
	theme := &Theme{
		Path:       themePath,
		Colors:     &SystemColors{},
		UserColors: make(UserColorGroup),
		Fonts:      NewFontCache(),
	}

	data, err := os.ReadFile(path.Join(themePath, "colors.json"))
	if err != nil {
		return nil, err
	}

	tc := map[string]string{}
	err = json.Unmarshal(data, &tc)

	for k, v := range tc {
		name := strings.ReplaceAll(k, ":", ".")
		theme.UserColors[name] = toColor(palette, v)
	}

	theme.Colors = Default.Colors

	if c := theme.UserColors.GetColor("background"); c != nil {
		theme.Colors.Background = c
	}

	if c := theme.UserColors.GetColor("foreground"); c != nil {
		theme.Colors.Foreground = c
	}

	if c := theme.UserColors.GetColor("text.primary"); c != nil {
		theme.Colors.TextPrimary = c
	}

	if c := theme.UserColors.GetColor("text.muted"); c != nil {
		theme.Colors.TextMuted = c
	}

	if c := theme.UserColors.GetColor("window.background"); c != nil {
		theme.Colors.WindowBackground = c
	}

	if c := theme.UserColors.GetColor("menu.background"); c != nil {
		theme.Colors.MenuBackground = c
	}

	if c := theme.UserColors.GetColor("border"); c != nil {
		theme.Colors.Border = c
	}

	data, err = os.ReadFile(path.Join(theme.Path, "fonts.json"))
	if err != nil {
		return nil, err
	}

	fontConfig := make(map[string]struct {
		Font string
		Size float64
	})

	if err := json.Unmarshal(data, &fontConfig); err != nil {
		return nil, err
	}

	for k, v := range fontConfig {
		font := path.Join(theme.Path, v.Font)
		if err := theme.Fonts.LoadFont(k, font, v.Size); err != nil {
			return nil, err
		}
	}

	return theme, nil
}
