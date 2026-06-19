package fonts

import (
	"io/ioutil"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

type Font struct {
	Width  int
	Height int
	Data   map[rune][]byte

	Face font.Face
}

func (f *Font) Measure(s string) (w, h int) {
	d := &font.Drawer{
		Face: f.Face,
	}
	a := d.MeasureString(s)

	return int(a >> 6), int(f.Face.Metrics().Height.Ceil())
}

func Load(path string, points float64) (*Font, error) {
	fontBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}
	face := truetype.NewFace(f, &truetype.Options{
		Size: points,
		// Hinting: font.HintingFull,
	})

	return &Font{
		Face:   face,
		Height: int(face.Metrics().Height.Ceil()),
	}, nil
}
