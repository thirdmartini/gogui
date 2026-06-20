package ux

import (
	"github.com/thirdmartini/gogui/pkg/ux/canvas"
)

type Drawable interface {
	//	Show(canvas canvas.Canvas)
	Draw(canvas canvas.Canvas)
	Visible(show bool)
}

type Widget interface {
	EventHandler
	Drawable
}

type ContainerProvider interface {
	AddWidget(name string, widget Widget) error
}

const (
	AlignDefault = 0x0
	AlignLeft    = 0x0
	AlignRight   = 0x1
	AlignCenter  = 0x2
	AlignTop     = 0x0
	AlignBottom  = 0x2
)

const (
	BorderLeft   = 0x1
	BorderRight  = 0x2
	BorderTop    = 0x4
	BorderBottom = 0x8
	BorderAll    = BorderLeft | BorderRight | BorderTop | BorderBottom
)
