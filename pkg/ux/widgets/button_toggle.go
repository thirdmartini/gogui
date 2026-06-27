package widgets

import (
	"image"

	"github.com/thirdmartini/gogui/pkg/log"

	"github.com/thirdmartini/gogui/pkg/ux"
	"github.com/thirdmartini/gogui/pkg/ux/canvas"
	"github.com/thirdmartini/gogui/pkg/ux/canvas/color"
	"github.com/thirdmartini/gogui/pkg/ux/themes"
)

type ToggleButton struct {
	*ux.Component

	colors []color.Color
	icons  []themes.Icon

	state uint

	OnClick func(state uint) bool
}

func (b *ToggleButton) Toggle(state uint) {
	if state >= uint(len(b.colors)) {
		state = 0
	}
	b.state = state
}

func (b *ToggleButton) Draw(canvas canvas.Canvas) {
	if !b.IsVisible() {
		return
	}
	x, y, w, h := b.X(), b.Y(), b.W(), b.H()
	if w < h {
		w = h
	}

	// for debugging
	//canvas.DrawRect(x, y, w, h, themes.GetColor("text.primary"), themes.GetColor("text.primary"))

	r := w / 2

	x = x + (w / 2)
	y = y + (h / 2)

	fillColor := b.colors[b.state]

	canvas.DrawCircle(x, y, r, 1, fillColor, fillColor)

	if b.state < uint(len(b.icons)) {
		// will dra centered
		b.icons[b.state].Draw(canvas, x, y, themes.GetColor("text.primary"))
	}

}

func (b *ToggleButton) OnEvent(event *ux.Event) bool {
	switch event.Type {
	case ux.EventTypeTouch:
		if b.InsidePoint(*event.Content.(*image.Point)) {
			log.Debugf("ToggleButton.OnEvent:%v in:%v", event, b.InsidePoint(*event.Content.(*image.Point)))
			b.state++
			b.state = b.state % uint(len(b.colors))
			return b.OnClick(b.state)
		}
	}
	return false
}

func NewToggleButton(name string, r image.Rectangle, colors []color.Color, icons []themes.Icon) *ToggleButton {
	return &ToggleButton{
		Component: ux.NewComponent(name, r),

		icons:  icons,
		colors: colors,

		OnClick: func(state uint) bool {
			log.Debugf("ToggleButton.OnClick:%v", state)
			return true
		},
	}
}
