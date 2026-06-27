package ux

import (
	"image"

	"github.com/thirdmartini/gogui/pkg/ux/canvas"
)

// Component is a shared class that contains common properties used by the draw code
type Component struct {
	name      string
	R         image.Rectangle
	isVisible bool
}

func (c *Component) Name() string {
	return c.name
}

func (c *Component) Inside(x, y int) bool {
	return c.R.Min.X <= x && x < c.R.Max.X && c.R.Min.Y <= y && y < c.R.Max.Y
}

func (c *Component) InsidePoint(point image.Point) bool {
	return c.Inside(point.X, point.Y)
}

func (c *Component) X() int {
	return c.R.Min.X
}

func (c *Component) Y() int {
	return c.R.Min.Y
}

func (c *Component) W() int {
	return c.R.Dx()
}

func (c *Component) H() int {
	return c.R.Dy()
}

func (c *Component) Rect() image.Rectangle {
	return c.R
}

func (c *Component) Visible(show bool) {
	c.isVisible = show
}

func (c *Component) IsVisible() bool {
	return c.isVisible
}

func (c *Component) Draw(_ canvas.Canvas) {
}

func (c *Component) OnEvent(_ *Event) bool {
	return false
}

func NewComponent(name string, r image.Rectangle) *Component {
	return &Component{
		name:      name,
		R:         r,
		isVisible: true,
	}
}

func NewComponentD(name string, x, y, w, h int) *Component {
	return &Component{
		name: name,
		R:    image.Rect(x, y, x+w, y+h),
	}
}
