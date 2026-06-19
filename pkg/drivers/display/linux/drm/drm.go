//go:build !linux

package drm

import (
	"fmt"
	"image"
)

type Display struct{}

func (d *Display) Draw(im *image.RGBA) error {
	return nil
}

func (d *Display) WithRotation(rotation int) *Display {
	return d
}

func (d *Display) Size() image.Point {
	return image.Point{}
}

func (d *Display) Close() error {
	return nil
}

func NewDisplay(device string) (*Display, error) {
	return nil, fmt.Errorf("not supported on this platform")
}
