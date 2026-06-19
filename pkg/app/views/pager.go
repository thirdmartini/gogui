package views

import (
	"image"

	"github.com/thirdmartini/gogui/pkg/ux"
	"github.com/thirdmartini/gogui/pkg/ux/themes"

	"github.com/thirdmartini/gogui/pkg/ux/canvas"
)

type Window interface {
	ux.Drawable
	ux.EventHandler
}

type Pager struct {
	views []Window
	idx   int

	background image.Image

	current Window
}

func (p *Pager) Next() Window {
	p.idx++
	if p.idx >= len(p.views) {
		p.idx = 0
	}
	p.current = p.views[p.idx]
	return p.current
}

func (p *Pager) Prev() Window {
	p.idx--
	if p.idx < 0 {
		p.idx = len(p.views) - 1
	}
	p.current = p.views[p.idx]
	return p.current
}

func (p *Pager) Select(idx int) {
	p.idx = idx
	if p.idx >= len(p.views) {
		p.idx = 0
	}
	p.current = p.views[p.idx]
}

func (p *Pager) Add(w Window) {
	p.views = append(p.views, w)
	if p.current == nil {
		p.current = w
	}
}

func (p *Pager) OnEvent(event *ux.Event) bool {
	switch event.Type {
	case ux.EventTypeButton, ux.EventTypeTouch:
		return p.current.OnEvent(event)

	case ux.EventTypeKey:
		switch uint8(event.Kind) {
		case ux.KeyPressUp: // Previous menu
			p.Prev()

		case ux.KeyPressDown: // Next Menu
			p.Next()

		default:
		}
	}
	return true
}

func (p *Pager) Draw(canvas canvas.Canvas) {
	if p.background != nil {
		canvas.DrawImage(0, 0, p.background)
	} else {
		canvas.Clear(themes.ColorBackground)
	}

	if p.current != nil {
		p.current.Draw(canvas)
	}
}

func (p *Pager) SetBackground(im image.Image) {
	p.background = im
}

func (p *Pager) Visible(show bool) {
}

func NewPager() *Pager {
	p := &Pager{}
	return p
}
