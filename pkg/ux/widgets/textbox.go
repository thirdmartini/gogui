package widgets

import (
	"image"

	"github.com/thirdmartini/gogui/pkg/ux"
	"github.com/thirdmartini/gogui/pkg/ux/canvas"
	"github.com/thirdmartini/gogui/pkg/ux/canvas/color"
	"github.com/thirdmartini/gogui/pkg/ux/canvas/fonts"
)

type TextBox struct {
	*ux.Component
	text      string
	fontColor color.Color
	bgColor   color.Color
	Font      *fonts.Font

	align      uint8
	borderOpts uint8
	visible    bool
}

func (t *TextBox) Align(align uint8) {
	t.align = align
}

func (t *TextBox) SetText(text string) {
	t.text = text
}

func (t *TextBox) Draw(canvas canvas.Canvas) {
	if t.visible {
		x1, y1 := t.X(), t.Y()
		fw, fh := t.Font.Measure(t.text)

		if t.bgColor != nil {
			canvas.DrawRect(x1, y1, fw+20, fh+10, t.bgColor, t.bgColor)
		}

		canvas.DrawText(x1+10, y1+fh+2, t.text, t.Font, t.fontColor)
	}
}

func (t *TextBox) Visible(show bool) {
	t.visible = show
}

func NewTextBox(name string, r image.Rectangle, align uint8, text string, font *fonts.Font, color color.Color, bg color.Color) *TextBox {
	return &TextBox{
		Component: ux.NewComponent(name, r),
		Font:      font,

		text:      text,
		align:     align,
		visible:   true,
		fontColor: color,
		bgColor:   bg,
	}
}
