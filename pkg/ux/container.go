package ux

import (
	"fmt"
)

// Container is a collection of widgets that directs user input
type Container struct {
	Widgets map[string]Widget
}

func (c *Container) OnEvent(event *Event) bool {
	for _, w := range c.Widgets {
		if w.OnEvent(event) {
			return true
		}
	}
	return false
}

func (c *Container) AddWidget(name string, widget Widget) error {
	if _, ok := c.Widgets[name]; ok {
		return fmt.Errorf("widget with name %s already exists", name)
	}
	c.Widgets[name] = widget
	return nil
}

func (c *Container) RemoveWidget(name string) {
	delete(c.Widgets, name)
}

func (c *Container) GetWidget(name string) (Widget, bool) {
	w, ok := c.Widgets[name]
	return w, ok
}

func NewContainer() *Container {
	return &Container{
		Widgets: make(map[string]Widget),
	}
}
