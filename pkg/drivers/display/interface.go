package display

import (
	"image"
)

const (
	RotationNone = 0
	Rotation90   = 1
	Rotation180  = 2
	Rotation270  = 3
)

type Display interface {
	Draw(im *image.RGBA) error
	Size() image.Point
	Close() error
}
