package ux

import (
	"image"

	"github.com/thirdmartini/gogui/pkg/drivers/display"
	"github.com/thirdmartini/gogui/pkg/ux/canvas"
)

type DisplayController struct {
	canvas.Canvas
	im      *image.RGBA
	view    Widget
	display display.Display
}

func (vc *DisplayController) OnEvent(event *Event) bool {
	return vc.view.OnEvent(event)
}

func (vc *DisplayController) OnRepaint() {
	vc.view.Draw(vc.Canvas)
	vc.display.Draw(vc.im)
}

func NewDisplayController(d display.Display, view Widget) *DisplayController {
	pt := d.Size()
	im := image.NewRGBA(image.Rect(0, 0, pt.X, pt.Y))

	vc := &DisplayController{
		Canvas:  canvas.NewGGCanvas(im),
		view:    view,
		im:      im,
		display: d,
	}
	vc.view.Visible(true)
	return vc
}
