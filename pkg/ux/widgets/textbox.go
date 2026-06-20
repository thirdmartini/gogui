package widgets

import (
	"fmt"

	"github.com/thirdmartini/gogui/pkg/ux"
	"github.com/thirdmartini/gogui/pkg/ux/canvas"
	"github.com/thirdmartini/gogui/pkg/ux/canvas/color"
	"github.com/thirdmartini/gogui/pkg/ux/canvas/fonts"
	"github.com/thirdmartini/gogui/pkg/ux/themes"
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
		canvas.DrawRect(t.x, t.y, t.w, t.h, themes.Default.Colors.Background, themes.Default.Colors.Background)

		if t.borderOpts&ux.BorderLeft == ux.BorderLeft {
			canvas.DrawLine(t.x, t.y, t.x, t.y+t.h, themes.Default.Colors.Border)
		}

		if t.borderOpts&ux.BorderRight == ux.BorderRight {
			canvas.DrawLine(t.x+t.w, t.y, t.x+t.w, t.y+t.h, themes.Default.Colors.Border)
		}

		if t.borderOpts&ux.BorderTop == ux.BorderTop {
			canvas.DrawLine(t.x, t.y, t.x+t.w, t.y, themes.Default.Colors.Border)
		}

		if t.borderOpts&ux.BorderBottom == ux.BorderBottom {
			canvas.DrawLine(t.x, t.y+t.h, t.x+t.w, t.y+t.h, themes.Default.Colors.Border)
		}

		fmt.Printf("TextBox.Draw: %v\n", t.text)
		canvas.DrawText(t.x, t.y+t.Font.Height, t.text, t.Font, themes.Default.Colors.TextPrimary)

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
