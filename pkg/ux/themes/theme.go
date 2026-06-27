package themes

import (
	"bytes"
	"encoding/json"
	"image"
	"log"
	"os"
	"path"

	"github.com/thirdmartini/gogui/pkg/ux/canvas/color"
)

type Theme struct {
	basePath    string
	searchPaths []string
	userColors  *UserColorGroup
	fonts       *FontCache
	palette     color.Palette
	icons       IconProvider
}

// NewColor creates a new color object (or returns one from the color cache)
// name assigns a human name to the color (like mybutton.background) and is optional
// hex is the html formated hex color 9IE  #FF0000
func (t *Theme) NewColor(name string, hex string) color.Color {
	return t.userColors.NewColor(t.palette, name, hex)
}

func (t *Theme) GetColor(name string) color.Color {
	return t.userColors.GetColor(name)
}

func (t *Theme) GetIcon(name string) Icon {
	return t.icons.GetIcon(name)
}

func (t *Theme) LoadImage(name string) image.Image {
	data, err := t.ReadFile(name)
	if err != nil {
		log.Printf("Warning: error opening image %s: %s", name, err)
		return nil
	}

	r := bytes.NewReader(data)

	// Decode will figure out what type of image is in the file on its own.
	// We just have to be sure all the image packages we want are imported.
	src, _, err := image.Decode(r)
	if err != nil {
		log.Printf("Warning: error opening image %s: %s", name, err)
		return nil
	}
	return src
}

func (t *Theme) ReadFile(name string) ([]byte, error) {
	filePath, err := t.FindFile(name)
	if err != nil {
		return nil, err
	}
	return os.ReadFile(filePath)
}

func (t *Theme) FindFile(name string) (string, error) {
	for _, s := range t.searchPaths {
		p := path.Join(t.basePath, s, name)
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", os.ErrNotExist
}

func Load(themePath string, palette color.Palette) (*Theme, error) {
	data, err := os.ReadFile(themePath)
	if err != nil {
		return nil, err
	}

	var themeConfig ThemeConfig
	err = json.Unmarshal(data, &themeConfig)
	if err != nil {
		return nil, err
	}

	theme := &Theme{
		basePath:    path.Dir(themePath),
		searchPaths: themeConfig.SearchOrder,
		palette:     palette,
		userColors:  NewColorsGroup(),
		fonts:       NewFontCache(),
	}

	for k, v := range themeConfig.Colors {
		theme.NewColor(k, v)
	}

	for k, v := range themeConfig.Fonts {
		fontPath, err := theme.FindFile(v.Font)
		if err != nil {
			return nil, err
		}
		if err := theme.fonts.LoadFont(k, fontPath, v.Size); err != nil {
			return nil, err
		}
	}

	if themeConfig.FontIcons != "" {
		filePath, err := theme.FindFile(themeConfig.FontIcons)
		if err != nil {
			return nil, err
		}

		theme.icons, err = NewIconFontProvider(filePath, 32)
		if err != nil {
			return nil, err
		}
	} else {
		// no icons provider
	}

	return theme, nil
}
