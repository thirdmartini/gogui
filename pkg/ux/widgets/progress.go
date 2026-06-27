package widgets

import (
	"image"

	"github.com/thirdmartini/gogui/pkg/ux"
	"github.com/thirdmartini/gogui/pkg/ux/canvas"
	"github.com/thirdmartini/gogui/pkg/ux/canvas/color"
	"github.com/thirdmartini/gogui/pkg/ux/themes"
)

type ProgressBar struct {
	*ux.Component
	progress float64
	fgColor  color.Color
	bgColor  color.Color
}

func (b *ProgressBar) SetProgress(progress float64) {
	b.progress = progress
}

func (b *ProgressBar) Draw(canvas canvas.Canvas) {
	if !b.IsVisible {
		return
	}

	x, y, w, h := b.X(), b.Y(), b.W(), b.H()

	r := h / 4
	if r > w/2 {
		r = w / 2
	}
	pw := int(float64(w) * b.progress / 100.0)

	// Thick green border: outer fill then inset white fill.
	canvas.DrawRoundedRect(x, y, w, h, r, b.bgColor, b.bgColor)
	canvas.DrawRoundedRect(x, y, pw, h, r, b.fgColor, b.fgColor)
}

func NewProgressBar(name string, r image.Rectangle) *ProgressBar {
	return &ProgressBar{
		Component: ux.NewComponent(name, r),
		fgColor:   themes.NewColor("progressbar.foreground", "#47A952"),
		bgColor:   themes.NewColor("progressbar.background", "#FFFFFF"),
		progress:  0,
	}
}
