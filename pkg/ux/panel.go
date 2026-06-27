package ux

import (
	"image"

	"github.com/thirdmartini/gogui/pkg/ux/canvas"
	"github.com/thirdmartini/gogui/pkg/ux/canvas/color"
)

type Panel struct {
	*Container
	*Component

	rect    image.Rectangle
	visible bool

	backgroundColor color.Color
}

func (p *Panel) Draw(canvas canvas.Canvas) {
	if !p.visible {
		return
	}

	x, y, w, h := p.X(), p.Y(), p.W(), p.H()

	r := 20
	border := max(3, h/28)

	borderColor := p.backgroundColor
	fillColor := p.backgroundColor

	// Thick green border: outer fill then inset white fill.
	canvas.DrawRoundedRect(x, y, w, h, r, borderColor, borderColor)
	innerR := r - border
	if innerR < 0 {
		innerR = 0
	}
	canvas.DrawRoundedRect(x+border, y+border, w-2*border, h-2*border, innerR, fillColor, fillColor)

	for _, widget := range p.Widgets {
		widget.Draw(canvas)
	}
}

func (p *Panel) Visible(show bool) {
	p.visible = show
}

func NewPanel(name string, rect image.Rectangle, bg color.Color) *Panel {
	p := &Panel{
		Container:       NewContainer(),
		Component:       NewComponent(name, rect),
		backgroundColor: bg,
		rect:            rect,
		visible:         true,
	}
	return p
}
