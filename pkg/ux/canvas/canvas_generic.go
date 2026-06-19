package canvas

import (
	"image"
	"math"

	"github.com/thirdmartini/gogui/pkg/ux/canvas/color"
	"github.com/thirdmartini/gogui/pkg/ux/canvas/fonts"
)

type Rect struct {
	x1 int
	y1 int
	x2 int
	y2 int
}

type GenericCanvas struct {
	surface Surface
	r       Rect // clipping region
	Repaint bool
}

func (im *GenericCanvas) setRect(x, y, w, h int) {
	im.r.x1 = x
	im.r.y1 = y
	im.r.x2 = x + w
	im.r.y2 = y + h

	if im.r.x2 > im.surface.Width() {
		im.r.x2 = im.surface.Width()
	}

	if im.r.y2 > im.surface.Height() {
		im.r.y2 = im.surface.Height()
	}
}

func (im *GenericCanvas) clearRect() {
	im.r.x1 = 0
	im.r.y1 = 0
	im.r.x2 = im.surface.Width()
	im.r.y2 = im.surface.Height()
}

func (im *GenericCanvas) DrawChar(x, y int, r rune, font *fonts.Font, fg, bg color.Color) {
	fontData, ok := font.Data[r]
	if !ok {
		return
	}

	ptr := 0
	for row := 0; row < font.Height; row++ {
		for col := 0; col < font.Width; col++ {
			if fontData[ptr]&(0x80>>(col%8)) != 0 {
				if fg != nil {
					im.DrawPixel(x+col, y+row, fg)
				}
			} else if bg != nil {
				im.DrawPixel(x+col, y+row, bg)
			}

			// font is wider then 8bits (8pixels).. so move to the next byte in the font data matrix
			if col%8 == 7 {
				ptr++
			}
		}
		// if we did not consume all the bits in the byte for width, move onto the t=next byte
		// since we're starting a new scan line
		if font.Width%8 != 0 {
			ptr++
		}
	}
}

func (im *GenericCanvas) DrawPixel(x, y int, c color.Color) {
	if (x >= im.r.x2) || (y >= im.r.y2) {
		return
	}
	im.surface.Set(x, y, c)
}

func (im *GenericCanvas) Clear(color color.Color) {
	for x := 0; x < im.surface.Width(); x++ {
		for y := 0; y < im.surface.Height(); y++ {
			im.DrawPixel(x, y, color)
		}
	}
}

func (im *GenericCanvas) DrawText(x, y int, text string, font *fonts.Font, fg color.Color) {
	xPoint := x
	yPoint := y

	for _, char := range text {
		im.DrawChar(xPoint, yPoint, char, font, fg, nil)
		xPoint += font.Width
	}
}

func (c *GenericCanvas) DrawTextWrapped(x, y, w, s int, text string, font *fonts.Font, fg color.Color) {

}

func (c *GenericCanvas) DrawImage(x, y int, image image.Image) {
}

func (c *GenericCanvas) ClipSet(x, y, w, h int) {
}

func (c *GenericCanvas) ClipReset() {
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// https://github.com/miloyip/line
func (im *GenericCanvas) DrawLine(x0, y0, x1, y1 int, c color.Color) {
	dx := abs(x1 - x0)
	sx := -1
	if x0 < x1 {
		sx = 1
	}
	dy := abs(y1 - y0)
	sy := -1
	if y0 < y1 {
		sy = 1
	}

	var err int
	if dx > dy {
		err = dx / 2
	} else {
		err = -dy / 2
	}

	for im.DrawPixel(x0, y0, c); x0 != x1 || y0 != y1; im.DrawPixel(x0, y0, c) {
		e2 := err
		if e2 > -dx {
			err -= dy
			x0 += sx
		}
		if e2 < dy {
			err += dx
			y0 += sy
		}
	}
}

func (im *GenericCanvas) DrawRect(x1, y1, w, h int, c color.Color, fill color.Color) {
	if fill != nil {
		for sx := x1 + 1; sx < x1+w; sx++ {
			for sy := y1 + 1; sy < y1+h; sy++ {
				im.DrawPixel(sx, sy, fill)
			}
		}
	}
	im.DrawLine(x1, y1, x1+w, y1, c)     // top
	im.DrawLine(x1, y1+h, x1+w, y1+h, c) // bottom

	im.DrawLine(x1, y1, x1, y1+h, c)     // left
	im.DrawLine(x1+w, y1, x1+w, y1+h, c) // right
}

func (im *GenericCanvas) DrawRoundedRect(x1, y1, w, h, r int, c color.Color, fill color.Color) {
	im.DrawRect(x1, y1, w, h, c, fill)
}

func (im *GenericCanvas) SetPixel(x, y int, color color.Color) {
	if x < 0 || y < 0 {
		return
	}
	im.surface.Set(x, y, color)
}

func (im *GenericCanvas) DrawPoint(x, y, r int, c color.Color) {
	im.DrawCircle(x, y, r, c, nil)
}

func (im *GenericCanvas) DrawCircle(x, y, r int, c color.Color, fill color.Color) {
	if r < 0 {
		return
	}
	x1, y1, err := -r, 0, 2-2*r
	for {
		im.SetPixel(x-x1, y+y1, c)
		im.SetPixel(x-y1, y-x1, c)
		im.SetPixel(x+x1, y-y1, c)
		im.SetPixel(x+y1, y+x1, c)
		r = err
		if r > x1 {
			x1++
			err += x1*2 + 1
		}
		if r <= y1 {
			y1++
			err += y1*2 + 1
		}
		if x1 >= 0 {
			break
		}
	}
}

func (im *GenericCanvas) DrawArc(x, y, r int, start, stop int, color color.Color) {
	for cur := start; cur <= stop; cur++ {
		x0 := x + int(float64(r)*math.Cos(float64(cur)*0.017453))
		y0 := y + int(float64(r)*math.Sin(float64(cur)*0.017453))
		im.SetPixel(x0, y0, color)
	}
}

func (im *GenericCanvas) Invalidate() {
	im.Repaint = true
}

func (im *GenericCanvas) Show() {
	im.surface.Show()
}

func (im *GenericCanvas) ColorPalette() color.Palette {
	return im.surface.ColorPalette()
}

func NewGenericCanvas(surface Surface) *GenericCanvas {
	c := &GenericCanvas{
		surface: surface,
	}
	c.clearRect()
	return c
}
