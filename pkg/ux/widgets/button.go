package widgets

import (
	"image"
	"log"

	"github.com/thirdmartini/gogui/pkg/ux"
	"github.com/thirdmartini/gogui/pkg/ux/canvas"
	"github.com/thirdmartini/gogui/pkg/ux/canvas/color"
	"github.com/thirdmartini/gogui/pkg/ux/canvas/fonts"
	"github.com/thirdmartini/gogui/pkg/ux/themes"
)

type Button struct {
	*ux.Component
	text    string
	Font    *fonts.Font
	OnClick func() bool

	BorderColor     color.Color
	BackgroundColor color.Color
	TextColor       color.Color
}

func (b *Button) Draw(canvas canvas.Canvas) {
	if !b.IsVisible {
		return
	}
	x, y, w, h := b.X(), b.Y(), b.W(), b.H()
	if w <= 0 || h <= 0 {
		return
	}

	r := h / 4
	if r > w/2 {
		r = w / 2
	}
	border := max(3, h/28)

	borderColor := b.BorderColor
	fillColor := b.BackgroundColor
	textColor := b.TextColor

	// Thick green border: outer fill then inset white fill.
	canvas.DrawRoundedRect(x, y, w, h, r, borderColor, borderColor)

	innerR := r - border
	if innerR < 0 {
		innerR = 0
	}
	canvas.DrawRoundedRect(x+border, y+border, w-2*border, h-2*border, innerR, fillColor, fillColor)

	if b.Font != nil && b.Font.Face != nil && b.text != "" {
		//_,  := b.Font.Measure(b.text)
		textY := y + (h / 2) - 4
		canvas.DrawTextCentered(x+w/2, textY, b.text, b.Font, textColor)
	}
}

func (b *Button) OnEvent(event *ux.Event) bool {
	log.Printf("Button.OnEvent: %v", event)
	switch event.Type {
	case ux.EventTypeTouch:
		if b.InsidePoint(*event.Content.(*image.Point)) && b.OnClick != nil {
			return b.OnClick()
		}
	}
	return false
}

func NewButton(name string, r image.Rectangle, text string) *Button {
	return &Button{
		Component:       ux.NewComponent(name, r),
		text:            text,
		BorderColor:     themes.NewColor("button.border", "#47A952"),
		BackgroundColor: themes.NewColor("button.background", "#FFFFFF"),
		TextColor:       themes.NewColor("button.text", "#000000"),
		Font:            themes.Font("default:medium"),

		OnClick: func() bool {
			log.Printf("Button.OnClick")
			return true
		},
	}
}
