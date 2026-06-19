package widgets

import (
	"image"
	"log"

	"github.com/thirdmartini/gogui/pkg/ux"
	"github.com/thirdmartini/gogui/pkg/ux/canvas/color"
	"github.com/thirdmartini/gogui/pkg/ux/canvas/fonts"
)

type Button struct {
	*TextBox
	OnClick func() bool
}

func (b *Button) OnEvent(event *ux.Event) bool {
	log.Printf("Button.OnEvent: %v", event)
	switch event.Type {
	case ux.EventTypeTouch:
		point := event.Content.(*image.Point)

		if point.X >= b.x && point.X <= b.x+b.w && point.Y >= b.y && point.Y <= b.y+b.h {
			if b.OnClick != nil {
				return b.OnClick()
			}
		}
	}
	return false
}

func NewButton(x, y, w, h int, align uint8, text string, font *fonts.Font, color color.Color) *Button {
	return &Button{
		TextBox: NewTextBox(x, y, w, h, align, text, font, color),
		OnClick: func() bool {
			log.Printf("Button.OnClick")
			return true
		},
	}
}
