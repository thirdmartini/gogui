package composer

import (
	"log"

	"github.com/thirdmartini/gogui/pkg/ux"
)

// Composer will construct a UI from a config file defining the compoonetns ath their relationships
type Composer struct {
	root       []string
	components map[string]interface{}
}

// Root will return the root widget of the given display number. We usually only have 1 display attached to our PC,
// but an RPI4/RPI5 can have dual displays, this lets us have a different ui composed on each display
func (c *Composer) Root(id uint) ux.Widget {
	if id >= uint(len(c.root)) {
		log.Panicf("invalid root id %d", id)
	}

	return c.GetWidget(c.root[id])
}

// Get will return the component with the given name the user is responsible for making sure they cast it to the
// correct type before using
func (c *Composer) Get(name string) interface{} {
	component, ok := c.components[name]
	if !ok {
		log.Panicf("no component with name " + name)
	}

	return component
}

// GetWidget will return the component with the given name cast to a Widget
func (c *Composer) GetWidget(name string) ux.Widget {
	w, ok := c.Get(name).(ux.Widget)
	if !ok {
		log.Panicf("component with name " + name + " is not a widget")
	}
	return w
}
