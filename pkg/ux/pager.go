package ux

import (
	"image"

	"github.com/thirdmartini/gogui/pkg/ux/canvas"
	"github.com/thirdmartini/gogui/pkg/ux/canvas/color"
	"github.com/thirdmartini/gogui/pkg/ux/themes"
)

var _ Container = (*Pager)(nil)

type Pager struct {
	*Component // for later shinanigans

	views      []Widget
	idx        int
	background image.Image
	current    Widget

	bgColor color.Color
}

func (p *Pager) Next() Widget {
	idx := p.idx + 1
	if idx >= len(p.views) {
		idx = 0
	}
	return p.Select(idx)
}

func (p *Pager) Prev() Widget {
	idx := p.idx - 1
	if idx < 0 {
		idx = len(p.views) - 1
	}
	return p.Select(idx)
}

func (p *Pager) Select(idx int) Widget {
	p.idx = idx
	if p.idx >= len(p.views) {
		p.idx = 0
	}

	p.current.Visible(false)
	p.current = p.views[p.idx]
	p.current.Visible(true)
	return p.current
}

func (p *Pager) AddWidget(name string, w Widget) error {
	p.Add(w)
	return nil
}

// fixme: deprecate this
func (p *Pager) Add(w Widget) {
	w.Visible(false)
	p.views = append(p.views, w)
	if p.current == nil {
		p.current = w
		w.Visible(true)
	}
}

func (p *Pager) OnEvent(event *Event) bool {
	switch event.Type {
	case ScreenSwipeLeft:
		_ = p.Next()
		return true

	case ScreenSwipeRight:
		_ = p.Prev()
		return true

	case EventTypeKey:
		switch uint8(event.Kind) {
		case KeyPressLeft: // Previous menu
			_ = p.Prev()
			return true

		case KeyPressRight: // Next Menu
			_ = p.Next()
			return true
		}
	}
	return p.current.OnEvent(event)
}

func (p *Pager) Draw(canvas canvas.Canvas) {
	// pager does not paint , it assumes children paint over it
	/*
		if p.background != nil {
			canvas.DrawImage(0, 0, p.background)
		} else {
			canvas.Clear(p.bgColor)
		}*/

	if p.current != nil {
		p.current.Draw(canvas)
	}
}

func (p *Pager) SetBackground(im image.Image) {
	p.background = im
}

func NewPager(name string, r image.Rectangle) *Pager {
	p := &Pager{
		Component: NewComponent(name, r),
		bgColor:   themes.NewColor("background", "#000000"),
	}
	return p
}
