package widgets

import (
	"github.com/thirdmartini/gogui/pkg/ux"
	"github.com/thirdmartini/gogui/pkg/ux/canvas"
	"github.com/thirdmartini/gogui/pkg/ux/canvas/color"
	"github.com/thirdmartini/gogui/pkg/ux/canvas/fonts"
)

type TextLabel struct {
	text       string
	Color      color.Color
	Background color.Color
	Font       *fonts.Font

	x       int
	y       int
	align   uint8
	visible bool
}

func (t *TextLabel) Align(align uint8) {
	t.align = align
}

func (t *TextLabel) SetText(text string) {
	t.text = text
}

func (t *TextLabel) SetColor(fg, bg color.Color) {
	t.Color = fg
	t.Background = bg
}

func (t *TextLabel) Draw(canvas canvas.Canvas) {
	if t.visible {
		canvas.DrawText(t.x, t.y, t.text, t.Font, t.Color)
	}
}

func (t *TextLabel) OnEvent(event *ux.Event) bool {
	return false
}

func (t *TextLabel) Visible(show bool) {
	t.visible = show
}

func NewTextLabel(x, y int, align uint8, text string, font *fonts.Font, fg, bg color.Color) *TextLabel {
	return &TextLabel{
		Font:       font,
		Color:      fg,
		Background: bg,

		x:       x,
		y:       y,
		text:    text,
		align:   align,
		visible: true,
	}
}
