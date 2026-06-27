package composer

import (
	"fmt"
	"log"

	"github.com/thirdmartini/gogui/pkg/ux"
	"github.com/thirdmartini/gogui/pkg/ux/themes"
	"github.com/thirdmartini/gogui/pkg/ux/widgets"
)

// common construct for builtin components
func (c *Composer) construct(def *ComponentDefinition) (interface{}, error) {

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
		widget = widgets.NewButton(def.Name, def.Rect(), def.Properties.Text)

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

	case "ux.widget.textbox":
		widget = widgets.NewTextBox(
			def.Name,
			def.Rect(),
			alignValue(def.Properties.Align),
			def.Properties.Text,
			themes.Font(def.Properties.Font),
			themes.GetColor(def.Properties.ColorText),
			themes.GetColor(def.Properties.ColorBackground),
		)

	case "ux.widget.progressbar":
		pb := widgets.NewProgressBar(
			def.Name,
			def.Rect(),
		)

		if v, ok := def.Custom["Progress"].(float64); ok {
			pb.SetProgress(v)
		}
		widget = pb

	case "ux.panel":
		widget = ux.NewPanel(def.Name, def.Rect(), themes.GetColor(def.Properties.ColorBackground))

	case "ux.pager":
		widget = ux.NewPager(def.Name, def.Rect())

	case "ux.window":
		w := ux.NewWindow(def.Rect())

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
