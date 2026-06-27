package composer

import (
	"fmt"

	"github.com/thirdmartini/gogui/pkg/ux"
)

type Constructor func(c *Composer, def *ComponentDefinition) (ux.Widget, error)

var constructors = make(map[string]Constructor)

func RegisterConstructor(name string, c Constructor) {
	constructors[name] = c
}
func GetConstructor(name string) (Constructor, error) {
	if c, ok := constructors[name]; ok {
		return c, nil
	}

	return nil, fmt.Errorf("no constructor for %s", name)
}
