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
	*BasicContainer
	background image.Image
}

func (w *Window) Draw(canvas canvas.Canvas) {
	if !w.IsVisible() {
		return
	}
	r := w.Rect()

	canvas.ClipSet(r.Min.X, r.Min.Y, r.Dx(), r.Dy())
	if w.background != nil {
		canvas.DrawImage(r.Min.X, r.Min.Y, w.background)
	} else {
		bgColor := themes.NewColor("backgeround", "#000000")
		canvas.DrawRect(r.Min.X, r.Min.Y, r.Dx(), r.Dx(), bgColor, bgColor)
	}

	w.BasicContainer.Draw(canvas)
	canvas.ClipReset()
}

func (w *Window) SetBackground(im image.Image) {
	w.background = im
}

func NewWindow(name string, rect image.Rectangle) *Window {
	w := &Window{
		BasicContainer: NewBasicContainer(name, rect),
	}
	return w
}
