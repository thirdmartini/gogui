package composer

import (
	"fmt"
	"image"
	"log"

	"github.com/thirdmartini/gogui/pkg/ux"
	"github.com/thirdmartini/gogui/pkg/ux/themes"
	"github.com/thirdmartini/gogui/pkg/ux/widgets"
)

// common construct for builtin components
func (c *Composer) construct(def *ComponentDefinition) (interface{}, error) {
	r := image.Rect(def.Properties.X, def.Properties.Y, def.Properties.Width, def.Properties.Height)

	var parent ux.ContainerProvider
	if def.Parent != "" {
		parent = c.Get(def.Parent).(ux.ContainerProvider)
		if parent == nil {
			return nil, fmt.Errorf("component[%s] wants parent[%s] which is not a container", def.Name, def.Parent)
		}
	}

	var widget ux.Widget

	// handle builtin type
	switch def.Type {
	case "ux.widget.button":
		widget = widgets.NewButton(
			def.Properties.X,
			def.Properties.Y,
			def.Properties.Width,
			def.Properties.Height,
			alignValue(def.Properties.Align),
			def.Properties.Text,
			themes.Font(def.Properties.Font),
			themes.GetColor(def.Properties.ColorText),
		)

	case "ux.widget.iconbutton":
		widget = widgets.NewIconButton(
			def.Properties.X,
			def.Properties.Y,
			def.Properties.Width,
			def.Properties.Height,
			themes.LoadImage(def.Properties.Icon),
		)

	case "ux.widget.textlabel":
		widget = widgets.NewTextLabel(
			def.Properties.X,
			def.Properties.Y,
			alignValue(def.Properties.Align),
			def.Properties.Text,
			themes.Font(def.Properties.Font),
			themes.GetColor(def.Properties.ColorText),
			themes.GetColor(def.Properties.ColorBackground),
		)

	case "ux.pager":
		widget = ux.NewPager(def.Name, r)

	case "ux.window":
		w := ux.NewWindow(r)

		if def.Properties.Background != "" {
			w.SetBackground(themes.LoadImage(def.Properties.Background))
		}
		widget = w
	default:
		construct, err := GetConstructor(def.Type)
		if err != nil {
			return nil, err
		}

		widget, err = construct(def)
		if err != nil {
			return nil, err
		}
	}

	if parent != nil {
		log.Printf("Adding widget %s to parent %s\n", def.Name, def.Parent)
		err := parent.AddWidget(def.Name, widget)
		if err != nil {
			log.Panicf("Error adding widget %s to parent %s: %s\n", def.Name, def.Parent, err)
		}

	} else {
		log.Printf("Adding widget %s to parent %s\n", def.Name, def.Parent)
	}

	return widget, nil
}

func (c *Composer) constructCustom(def ComponentDefinition) ux.Widget {
	panic(fmt.Errorf("unknown component type: %s", def.Type))
}
