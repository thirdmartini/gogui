package app

import (
	"fmt"
	"image"

	"github.com/thirdmartini/gogui/pkg/app/views"
	"github.com/thirdmartini/gogui/pkg/ux/themes"

	"github.com/thirdmartini/gogui/pkg/ux"
	"github.com/thirdmartini/gogui/pkg/ux/canvas"
)

type ViewPort interface {
	canvas.Canvas
	Show()
}

type Window interface {
	ux.Drawable
	ux.EventHandler
}

type Controller struct {
	vps    []ViewPort
	pagers []*views.Pager
	idx    int
}

func (c *Controller) Next() int {
	c.idx++
	if c.idx >= len(c.vps) {
		c.idx = 0
	}
	return c.idx
}

func (c *Controller) Prev() int {
	c.idx--
	if c.idx <= 0 {
		c.idx = len(c.vps) - 1
	}
	return c.idx
}

func (c *Controller) AddViewPort(v ViewPort, p *views.Pager) {
	p.Select(len(c.pagers))
	c.vps = append(c.vps, v)
	c.pagers = append(c.pagers, p)
	c.idx = len(c.vps) - 1
}

func (c *Controller) OnRepaint() {
	for idx, vp := range c.vps {
		c.pagers[idx].Draw(vp)
		if c.idx == idx {
			font := themes.Font(themes.FontHeader)
			vp.DrawRoundedRect(300, -10, 200, 40, 5, themes.ColorBorder, themes.ColorMenuBackground)
			vp.DrawRoundedRect(300, -10, 200, 40, 5, themes.ColorBorder, nil)

			display := fmt.Sprintf("DISPLAY %d\n", idx)
			vp.DrawText(350, 20, display, font, themes.ColorTextPrimary)

		}
		vp.Show()
	}
}

func (c *Controller) OnEvent(event *ux.Event) bool {
	// return true to always cause a repaint
	fmt.Printf("Event: %+v\n", event)

	switch event.Type {
	case ux.EventTypeTouch:
		fmt.Printf("Event: %d on %d  %+v\n", event.Kind, c.idx, event.Content.(*image.Point))

	case ux.EventTypeKey:
		fmt.Printf("Event: %d on %d\n", event.Kind, c.idx)
		switch uint8(event.Kind) {
		case ux.KeyPressUp, ux.KeyPressDown:
			if c.idx < len(c.vps) {
				c.pagers[c.idx].OnEvent(event)
			}
			return true

		case ux.KeyPressRight:
			c.Next()
			return true

		case ux.KeyPressLeft:
			c.Prev()
			return true

		default:
		}
	}

	return false
}
