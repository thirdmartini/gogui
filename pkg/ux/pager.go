package ux

import (
	"image"

	"github.com/thirdmartini/gogui/pkg/log"
	"github.com/thirdmartini/gogui/pkg/ux/canvas"
	"github.com/thirdmartini/gogui/pkg/ux/canvas/color"
	"github.com/thirdmartini/gogui/pkg/ux/themes"
)

var _ Container = (*Pager)(nil)

const (
	FlowDirectionHorizontal = iota
	FlowDirectionVertical
)

type Pager struct {
	*Component // for later shinanigans

	views      []Widget
	idx        int
	background image.Image
	current    Widget

	direction int
	wrap      bool

	bgColor color.Color
}

func (p *Pager) Next() Widget {
	idx := p.idx + 1
	if idx >= len(p.views) {
		if p.wrap {
			idx = 0
		} else {
			log.Debugf("pager: no more pages")
			idx = len(p.views) - 1
		}

	}
	return p.Select(idx)
}

func (p *Pager) Prev() Widget {
	idx := p.idx - 1
	if idx < 0 {
		if p.wrap {
			idx = len(p.views) - 1
		} else {
			log.Debugf("pager: no more pages")
			idx = 0
		}

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

const (
	ActionNone = iota
	ActionNextHorizontal
	ActionPrevHorizontal
	ActionNextVertical
	ActionPrevVertical
)

func (p *Pager) collapseAction(event *Event) int {
	log.Debugf("collapseAction: %v", event)
	switch event.Type {
	case ScreenSwipeLeft:
		return ActionNextHorizontal

	case ScreenSwipeRight:
		return ActionPrevHorizontal

	case ScreenSwipeUp:
		return ActionNextVertical

	case ScreenSwipeDown:
		return ActionPrevVertical

	case EventTypeKey:
		switch uint8(event.Kind) {
		case KeyPressLeft: // Previous menu
			return ActionPrevHorizontal

		case KeyPressRight: // Next Menu
			return ActionNextHorizontal

		case KeyPressUp: // Previous menu
			return ActionPrevVertical

		case KeyPressDown: // Next Menu
			return ActionNextVertical
		}
	}
	return 0
}

func (p *Pager) OnEvent(event *Event) bool {
	action := p.collapseAction(event)
	if action != 0 {
		switch p.direction {
		case FlowDirectionHorizontal:
			switch action {
			case ActionNextHorizontal:
				_ = p.Next()
				return true
			case ActionPrevHorizontal:
				_ = p.Prev()
				return true
			}

		case FlowDirectionVertical:
			switch action {
			case ActionNextVertical:
				_ = p.Next()
				return true
			case ActionPrevVertical:
				_ = p.Prev()
				return true
			}
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

func (p *Pager) SetFlowDirection(direction int) {
	p.direction = direction
}

func (p *Pager) SetWrap(wrap bool) {
	p.wrap = wrap
}

func NewPager(name string, r image.Rectangle) *Pager {
	p := &Pager{
		Component: NewComponent(name, r),
		bgColor:   themes.NewColor("background", "#000000"),
		direction: FlowDirectionHorizontal,
		wrap:      true,
	}
	return p
}
