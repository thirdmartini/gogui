package widgets

import (
	"image"

	"github.com/thirdmartini/gogui/pkg/ux/canvas"
	"github.com/thirdmartini/gogui/pkg/ux/canvas/color"
	"github.com/thirdmartini/gogui/pkg/ux/canvas/fonts"
	"github.com/thirdmartini/gogui/pkg/ux/themes"
)

type ScrollingTextBox struct {
	*TextBox
	text   string
	offset uint64
}

func (t *ScrollingTextBox) SetText(text string) {
	t.text = text
	t.offset = 0
}

func (t *ScrollingTextBox) Draw(canvas canvas.Canvas) {
	if len(t.text) == 0 {
		t.TextBox.Draw(canvas)
		return
	}

	idx := t.offset % uint64(len(t.text))

	count := uint64(t.TextBox.W()) / uint64(t.TextBox.Font.Width)

	t.TextBox.SetText(t.text[idx:])
	t.TextBox.Draw(canvas)
	t.offset++

	if idx+count > uint64(len(t.text)) {
		t.offset = 0
	}
}

func NewScrollingTextBox(name string, r image.Rectangle, align uint8, text string, font *fonts.Font, color color.Color) *ScrollingTextBox {
	return &ScrollingTextBox{
		TextBox: NewTextBox(name, r, align, text, font, color, themes.NewColor("textbox.background", "#FFFFFF")),
		text:    text,
	}
}
