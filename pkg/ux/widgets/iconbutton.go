package widgets

import (
	"image"
	"log"

	"github.com/thirdmartini/gogui/pkg/ux"
	"github.com/thirdmartini/gogui/pkg/ux/canvas"
	"github.com/thirdmartini/gogui/pkg/ux/canvas/color"
	"github.com/thirdmartini/gogui/pkg/ux/themes"
)

type IconButton struct {
	x, y, w, h int
	icon       image.Image

	visible bool

	BorderColor     color.Color
	BackgroundColor color.Color
	PressedColor    color.Color
	HighlightColor  color.Color
	ShadowColor     color.Color
	CornerRadius    int

	OnTouch func(pressed bool) bool
}

func (b *IconButton) contains(px, py int) bool {
	return px >= b.x && px < b.x+b.w && py >= b.y && py < b.y+b.h
}

func (b *IconButton) OnEvent(event *ux.Event) bool {
	switch event.Type {
	case ux.EventTypeTouch:
		point := event.Content.(*image.Point)
		if b.contains(point.X, point.Y) {
			if b.OnTouch != nil {
				return b.OnTouch(true)
			}
			return true
		}
	}
	return false
}

func (b *IconButton) Draw(canvas canvas.Canvas) {
	if !b.visible {
		return
	}

	x, y := b.x, b.y
	fill := b.BackgroundColor

	canvas.DrawRoundedRect(x, y, b.w, b.h, b.CornerRadius, b.BorderColor, fill)

	/*
		if b.pressed {
			canvas.DrawLine(x, y, x+b.w, y, b.ShadowColor)
			canvas.DrawLine(x, y, x, y+b.h, b.ShadowColor)
		} else {
			canvas.DrawLine(x, y, x+b.w, y, b.HighlightColor)
			canvas.DrawLine(x, y, x, y+b.h, b.HighlightColor)
			canvas.DrawLine(x, y+b.h, x+b.w, y+b.h, b.ShadowColor)
			canvas.DrawLine(x+b.w, y, x+b.w, y+b.h, b.ShadowColor)
		}*/

	if b.icon != nil {
		bounds := b.icon.Bounds()
		iconW := bounds.Dx()
		iconH := bounds.Dy()
		ix := x + (b.w-iconW)/2
		iy := y + (b.h-iconH)/2
		canvas.DrawImage(ix, iy, b.icon)
	}
}

func (b *IconButton) Visible(show bool) {
	b.visible = show
}

func NewIconButton(x, y, w, h int, icon image.Image) *IconButton {
	return &IconButton{
		x:       x,
		y:       y,
		w:       w,
		h:       h,
		icon:    icon,
		visible: true,

		BorderColor:     themes.NewColor("iconbutton.border", "#3A7BBD"),
		BackgroundColor: themes.NewColor("iconbutton.background", "#5B9BD5"),
		HighlightColor:  themes.NewColor("iconbutton.highlight", "#82B4E6"),
		ShadowColor:     themes.NewColor("iconbutton.shadow", "#234B78"),
		CornerRadius:    min(w, h) / 5,

		OnTouch: func(pressed bool) bool {
			log.Printf("IconButton.OnTouch: pressed=%v", pressed)
			return true
		},
	}
}
