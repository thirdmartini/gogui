package ux

import (
	"image"

	"github.com/thirdmartini/gogui/pkg/ux/canvas"
	"github.com/thirdmartini/gogui/pkg/ux/themes"
)

type PhysicalRect struct {
	x, y, w, h int
}

type Window struct {
	*Container

	rect    image.Rectangle
	visible bool

	background image.Image
}

func (w *Window) Draw(canvas canvas.Canvas) {
	if !w.visible {
		return
	}
	r := w.rect

	canvas.ClipSet(r.Min.X, r.Min.Y, r.Dx(), r.Dy())
	if w.background != nil {
		canvas.DrawImage(r.Min.X, r.Min.Y, w.background)
	} else {
		bgColor := themes.NewColor("backgeround", "#000000")
		canvas.DrawRect(r.Min.X, r.Min.Y, r.Dx(), r.Dx(), bgColor, bgColor)
	}

	for _, widget := range w.Widgets {
		widget.Draw(canvas)
	}
	canvas.ClipReset()
}

func (c *Window) OnEvent(event *Event) bool {
	for _, w := range c.Widgets {
		if w.OnEvent(event) {
			return true
		}
	}
	return false
}

func (w *Window) Visible(show bool) {
	w.visible = show
}

func (w *Window) SetBackground(im image.Image) {
	w.background = im
}

func NewWindow(rect image.Rectangle) *Window {
	w := &Window{
		Container: NewContainer(),
		rect:      rect,
	}
	return w
}
