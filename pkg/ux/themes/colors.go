package themes

import (
	"fmt"
	"log"
	"strings"

	"github.com/thirdmartini/gogui/pkg/ux/canvas/color"
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

func toColor(p color.Palette, hexColor string) color.Color {
	var r, g, b uint8

	if cnt, err := fmt.Sscanf(hexColor, "#%02x%02x%02x", &r, &g, &b); err == nil && cnt == 3 {
		return p.NewRGB8(r, g, b)
	}

	return p.NewRGB8(0, 0, 0)
}

func mustMakeColor(p color.Palette, hexColor string) color.Color {
	c := toColor(p, hexColor)
	if c == nil {
		log.Fatalf("no color [%s]", hexColor)
	}
	return c
}

type UserColorGroup struct {
	byName map[string]color.Color
	byHex  map[string]color.Color
}

func (ucg *UserColorGroup) GetColor(name string) color.Color {
	if strings.HasPrefix(name, "#") {
		if c, ok := ucg.byHex[name]; ok {
			return c
		}
	}

	if c, ok := ucg.byName[name]; ok {
		return c
	}
	log.Panicf("no color named [%s]", name)
	return nil
}

func (ucg *UserColorGroup) NewColor(palette color.Palette, name string, hex string) color.Color {
	if name != "" {
		if c, ok := ucg.byName[name]; ok {
			return c
		}
	} else {
		// if the name is blank look it up by hex color
		if c, ok := ucg.byHex[hex]; ok {
			return c
		}
	}

	c := mustMakeColor(palette, hex)
	log.Printf("New color %s/%s\n", name, hex)
	ucg.byHex[hex] = c
	ucg.byName[name] = c
	return c
}

func NewColorsGroup() *UserColorGroup {
	return &UserColorGroup{
		byName: make(map[string]color.Color),
		byHex:  make(map[string]color.Color),
	}
}
