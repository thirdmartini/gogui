package themes

import (
	"github.com/thirdmartini/gogui/pkg/ux/canvas/fonts"
)

var DefaultCache = NewFontCache()

type FontCache struct {
	fonts map[string]*fonts.Font
}

func (fc *FontCache) LoadFont(name string, path string, points float64) error {
	font, err := fonts.Load(path, points)
	if err != nil {
		return err
	}

	fc.fonts[name] = font
	return nil
}

func (fc *FontCache) Font(name string) *fonts.Font {
	if font, ok := fc.fonts[name]; ok {
		return font
	}
	return &fonts.Font{}
}

func NewFontCache() *FontCache {
	return &FontCache{
		fonts: make(map[string]*fonts.Font),
	}
}
