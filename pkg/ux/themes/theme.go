package themes

import (
	"encoding/json"
	"image"
	"log"
	"os"
	"path"

	"github.com/thirdmartini/gogui/pkg/ux/canvas/color"
)

type Theme struct {
	path       string
	userColors *UserColorGroup
	fonts      *FontCache
	palette    color.Palette
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

func (t *Theme) LoadImage(name string) image.Image {
	imgPath := path.Join(t.path, name)
	infile, err := os.Open(imgPath)
	if err != nil {
		log.Printf("Warning: error opening image %s: %s", imgPath, err)
		return nil
	}
	defer infile.Close()

	// Decode will figure out what type of image is in the file on its own.
	// We just have to be sure all the image packages we want are imported.
	src, _, err := image.Decode(infile)
	if err != nil {
		log.Printf("Warning: error opening image %s: %s", imgPath, err)
		return nil
	}
	return src
}

func Load(themePath string, palette color.Palette) (*Theme, error) {
	theme := &Theme{
		path:       themePath,
		palette:    palette,
		userColors: NewColorsGroup(),
		fonts:      NewFontCache(),
	}

	data, err := os.ReadFile(path.Join(themePath, "colors.json"))
	if err != nil {
		return nil, err
	}

	tc := map[string]string{}
	err = json.Unmarshal(data, &tc)

	for k, v := range tc {
		theme.NewColor(k, v)
	}

	data, err = os.ReadFile(path.Join(theme.path, "fonts.json"))
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
		font := path.Join(theme.path, v.Font)
		if err := theme.fonts.LoadFont(k, font, v.Size); err != nil {
			return nil, err
		}
	}

	return theme, nil
}
