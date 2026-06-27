package composer

import (
	"fmt"

	"github.com/thirdmartini/gogui/pkg/log"

	"github.com/thirdmartini/gogui/pkg/ux"
	"github.com/thirdmartini/gogui/pkg/ux/themes"
	"github.com/thirdmartini/gogui/pkg/ux/widgets"
)

// common construct for builtin components
func (c *Composer) construct(def *ComponentDefinition) (interface{}, error) {

	var parent ux.Container
	if def.Parent != "" {
		parent = c.Get(def.Parent).(ux.Container)
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

	case "ux.widget.togglebutton":
		widget = widgets.NewToggleButton(
			def.Name,
			def.Rect(),
			def.Colors(),
			def.Icons(),
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
		pager := ux.NewPager(def.Name, def.Rect())
		if def.Properties.Flow != "" {
			switch def.Properties.Flow {
			case "horizontal":
				pager.SetFlowDirection(ux.FlowDirectionHorizontal)
				pager.SetWrap(true)
			case "horizontal;nowrap":
				pager.SetFlowDirection(ux.FlowDirectionHorizontal)
				pager.SetWrap(false)
			case "vertical":
				pager.SetFlowDirection(ux.FlowDirectionVertical)
				pager.SetWrap(true)
			case "vertical;nowrap":
				pager.SetFlowDirection(ux.FlowDirectionVertical)
				pager.SetWrap(false)
			}
		}
		widget = pager

	case "ux.window":
		w := ux.NewWindow(def.Name, def.Rect())

		if def.Properties.Background != "" {
			w.SetBackground(themes.LoadImage(def.Properties.Background))
		}
		widget = w
	default:
		construct, err := GetConstructor(def.Type)
		if err != nil {
			return nil, err
		}

		widget, err = construct(c, def)
		if err != nil {
			return nil, err
		}
	}

	if parent != nil {
		err := parent.AddWidget(def.Name, widget)
		if err != nil {
			log.Panicf("Error adding widget %s to parent %s: %s", def.Name, def.Parent, err)
		}
		log.Debugf("Added widget %s to parent %s", def.WidgetString(), def.Parent)
	} else {
		log.Debugf("Created widget %s", def.WidgetString())
	}

	return widget, nil
}

func (c *Composer) constructCustom(def ComponentDefinition) ux.Widget {
	panic(fmt.Errorf("unknown component type: %s", def.WidgetString()))
}
