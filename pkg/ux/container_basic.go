package ux

import (
	"fmt"
	"image"

	"github.com/thirdmartini/gogui/pkg/ux/canvas"
)

// BasicContainer is a collection of widgets that directs user input
type BasicContainer struct {
	*Component
	Widgets map[string]Widget
}

func (c *BasicContainer) AddWidget(name string, widget Widget) error {
	if _, ok := c.Widgets[name]; ok {
		return fmt.Errorf("widget with name %s already exists", name)
	}
	c.Widgets[name] = widget
	return nil
}

func (c *BasicContainer) RemoveWidget(name string) {
	delete(c.Widgets, name)
}

func (c *BasicContainer) GetWidget(name string) (Widget, bool) {
	w, ok := c.Widgets[name]
	return w, ok
}

func (c *BasicContainer) Draw(canvas canvas.Canvas) {
	if !c.isVisible {
		return
	}
	for _, w := range c.Widgets {
		w.Draw(canvas)
	}
}

func (c *BasicContainer) OnEvent(event *Event) bool {
	for _, w := range c.Widgets {
		if w.OnEvent(event) {
			return true
		}
	}
	return false
}

func NewBasicContainer(name string, r image.Rectangle) *BasicContainer {
	return &BasicContainer{
		Component: NewComponent(name, r),
		Widgets:   make(map[string]Widget),
	}
}
