package composer

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/thirdmartini/gogui/pkg/ux"
)

// ComponentDefinition represents the structure defining a UI component's type, name, parent, and associated properties.
type ComponentDefinition struct {
	Type       string
	Name       string
	Parent     string
	Properties struct { // These are commonly used properties by components
		Text            string
		X               int
		Y               int
		Width           int
		Height          int
		Align           string
		Font            string
		ColorText       string
		ColorBackground string
		Background      string
		Icon            string
	}
	Custom map[string]interface{} // Custom properties for use by external components can be put here
}

// ComposerConfig defines the configuration for setting up the UI composer, including theme, displays, and components.
type ComposerConfig struct {
	Theme      string // not supported yet
	Displays   []string
	Components []ComponentDefinition
}

// From construct an application from a config file
func From(configPath string) (*ux.Application, *Composer, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, nil, err
	}

	config := ComposerConfig{}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, nil, err
	}

	app := ux.NewApplication()
	if config.Theme != "" {
		app.WithTheme(config.Theme)
	}

	composer := &Composer{
		root:       config.Displays,
		components: make(map[string]interface{}),
	}

	for idx := range config.Components {
		def := config.Components[idx]

		if _, ok := composer.components[def.Name]; ok {
			return nil, nil, fmt.Errorf("duplicate component name %s", def.Name)
		}

		component, err := composer.construct(&def)
		if err != nil {
			return nil, nil, err
		}
		composer.components[def.Name] = component
	}

	return app, composer, nil
}
