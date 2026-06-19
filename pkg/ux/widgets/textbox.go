package widgets

import (
	"fmt"

	"github.com/thirdmartini/gogui/pkg/ux/canvas"
	"github.com/thirdmartini/gogui/pkg/ux/canvas/color"
	"github.com/thirdmartini/gogui/pkg/ux/canvas/fonts"
)

const (
	BorderLeft   = 0x1
	BorderRight  = 0x2
	BorderTop    = 0x4
	BorderBottom = 0x8
	BorderAll    = BorderLeft | BorderRight | BorderTop | BorderBottom
)

type TextBox struct {
	text            string
	FontColor       color.Color
	BorderColor     color.Color
	BackgroundColor color.Color
	Font            *fonts.Font

	x          int
	y          int
	w          int
	h          int
	align      uint8
	borderOpts uint8
	visible    bool
}

func (t *TextBox) Width() int {
	return t.w
}

func (t *TextBox) Align(align uint8) {
	t.align = align
}

func (t *TextBox) SetText(text string) {
	t.text = text
}

func (t *TextBox) SetFontColor(fg color.Color) {
	t.FontColor = fg
}

func (t *TextBox) SetBorder(border uint8, borderColor color.Color) {
	t.borderOpts = border
	t.BorderColor = borderColor
}

func (t *TextBox) SetBackground(bgColor color.Color) {
	t.BackgroundColor = bgColor
}

func (t *TextBox) Draw(canvas canvas.Canvas) {
	if t.visible {
		fmt.Printf(" + TextBox:Draw() %+v\n", t.BackgroundColor)
		canvas.DrawRect(t.x, t.y, t.w, t.h, t.BackgroundColor, t.BackgroundColor)

		if t.borderOpts&BorderLeft == BorderLeft {
			fmt.Printf(" + BorderColor %+v\n", t.BorderColor)
			canvas.DrawLine(t.x, t.y, t.x, t.y+t.h, t.BorderColor)
		}

		if t.borderOpts&BorderRight == BorderRight {
			canvas.DrawLine(t.x+t.w, t.y, t.x+t.w, t.y+t.h, t.BorderColor)
		}

		if t.borderOpts&BorderTop == BorderTop {
			canvas.DrawLine(t.x, t.y, t.x+t.w, t.y, t.BorderColor)
		}

		if t.borderOpts&BorderBottom == BorderBottom {
			canvas.DrawLine(t.x, t.y+t.h, t.x+t.w, t.y+t.h, t.BorderColor)
		}

		fmt.Printf(" + FontColor %+v %s %v/%v\n", t.FontColor, t.text, t.Font.Width, t.Font.Height)

		canvas.DrawText(t.x, t.y+t.Font.Height, t.text, t.Font, t.FontColor)

	}
}

func (t *TextBox) Visible(show bool) {
	t.visible = show
}

func NewTextBox(x, y, w, h int, align uint8, text string, font *fonts.Font, color color.Color) *TextBox {
	return &TextBox{
		Font:      font,
		FontColor: color,
		x:         x,
		y:         y,
		w:         w,
		h:         h,
		text:      text,
		align:     align,
		visible:   true,
	}
}
