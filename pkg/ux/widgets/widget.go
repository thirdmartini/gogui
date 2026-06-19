package widgets

import (
	"github.com/thirdmartini/gogui/pkg/ux"
)

type Widget interface {
	ux.Drawable
	ux.EventHandler
}
