package canvas

import (
	"image"

	"github.com/fogleman/gg"

	"github.com/thirdmartini/gogui/pkg/ux/canvas/color"
	"github.com/thirdmartini/gogui/pkg/ux/canvas/fonts"
)

type CanvasGG struct {
	gg *gg.Context
}

func (c *CanvasGG) setColor(color color.Color) {
	r, g, b := color.RGB8()
	c.gg.SetRGBA255(int(r), int(g), int(b), 255)
}

func (c *CanvasGG) Clear(color color.Color) {
	//fmt.Printf("Clear: %+v\n", color)
	c.setColor(color)
	c.gg.Clear()
}

func (c *CanvasGG) DrawPixel(x, y int, color color.Color) {
	c.setColor(color)
	c.gg.DrawPoint(float64(x), float64(y), 1)
}

func (c *CanvasGG) DrawPoint(x, y, r int, color color.Color) {
	c.setColor(color)
	c.gg.DrawPoint(float64(x), float64(y), float64(r))
}

func (c *CanvasGG) DrawTextCentered(x, y int, text string, font *fonts.Font, fg color.Color) {
	if font.Face != nil {
		fw, fh := font.Measure(text)

		fx := x - (fw / 2)
		fy := y + (fh / 2)

		c.setColor(fg)
		c.gg.SetFontFace(font.Face)
		c.gg.DrawString(text, float64(fx), float64(fy))
	}
}

func (c *CanvasGG) DrawText(x, y int, text string, font *fonts.Font, fg color.Color) {
	if font.Face != nil {
		c.setColor(fg)
		c.gg.SetFontFace(font.Face)
		c.gg.DrawString(text, float64(x), float64(y))
	}
}

func (c *CanvasGG) DrawTextWrapped(x, y, w, s int, text string, font *fonts.Font, fg color.Color) {
	if font.Face != nil {
		c.setColor(fg)
		c.gg.SetFontFace(font.Face)
		c.gg.DrawStringWrapped(text, float64(x), float64(y), 0, 0, float64(w), float64(s), gg.AlignLeft)
	}
}

func (c *CanvasGG) DrawLine(x1, y1, x2, y2 int, color color.Color) {
	c.setColor(color)

	c.gg.SetLineWidth(1)
	c.gg.DrawLine(float64(x1)+.5,
		float64(y1)+.5,
		float64(x2)+.5,
		float64(y2)+.5)

	c.gg.Stroke()
}

func (c *CanvasGG) DrawRect(x1, y1, w, h int, frame color.Color, fill color.Color) {
	c.setColor(frame)
	c.gg.DrawRectangle(float64(x1),
		float64(y1),
		float64(w),
		float64(h))
	if fill != nil {
		c.gg.Fill()
	}
	c.gg.Stroke()
}

func (c *CanvasGG) DrawRoundedRect(x1, y1, w, h, r int, frame color.Color, fill color.Color) {
	c.setColor(frame)
	c.gg.DrawRoundedRectangle(float64(x1),
		float64(y1),
		float64(w),
		float64(h),
		float64(r))
	if fill != nil {
		c.setColor(fill)
		c.gg.Fill()
	}
	c.gg.Stroke()
}

func (c *CanvasGG) DrawImage(x, y int, im image.Image) {
	c.gg.DrawImage(im, x, y)
}

func (c *CanvasGG) ClipSet(x, y, w, h int) {
	c.gg.DrawRectangle(float64(x), float64(y), float64(w), float64(h))
	c.gg.Clip()
}

func (c *CanvasGG) ClipReset() {
	c.gg.ResetClip()
}

func (c *CanvasGG) DrawCircle(x, y, r, w int, color color.Color, fill color.Color) {
	c.setColor(color)
	c.gg.SetLineWidth(float64(w))
	c.gg.DrawCircle(float64(x), float64(y), float64(r))
	if fill != nil {
		c.setColor(fill)
		c.gg.Fill()
	}
	c.gg.Stroke()
	c.gg.SetLineWidth(1)
}

func (c *CanvasGG) DrawArc(x, y, r, w int, start, stop int, color color.Color, fill color.Color) {
	c.setColor(color)
	c.gg.SetLineWidth(float64(w))
	c.gg.DrawArc(float64(x), float64(y), float64(r), gg.Radians(float64(start)), gg.Radians(float64(stop)))
	if fill != nil {
		c.setColor(fill)
		c.gg.Fill()
	}
	c.gg.Stroke()
	c.gg.SetLineWidth(1)
}

func (c *CanvasGG) DrawEllipticalArc(x, y, rx, ry int, start, stop int, color color.Color, fill color.Color) {
	c.setColor(color)
	c.gg.DrawEllipticalArc(float64(x), float64(y), float64(rx), float64(ry), gg.Radians(float64(start)), gg.Radians(float64(stop)))
	if fill != nil {
		c.setColor(fill)
		c.gg.Fill()
	}
	c.gg.Stroke()
}

func (c *CanvasGG) Show() {
}

func (c *CanvasGG) ColorPalette() color.Palette {
	return &color.Palette888{}
}

func (c *CanvasGG) Invalidate() {

}

func NewGGCanvas(im *image.RGBA) *CanvasGG {
	g := gg.NewContextForRGBA(im)
	g.SetLineCap(gg.LineCapButt)

	return &CanvasGG{
		gg: g,
	}
}
