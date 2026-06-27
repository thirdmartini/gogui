package ux

import (
	"image"
)

const (
	ButtonNone  = 0
	ButtonOne   = 1
	ButtonTwo   = 2
	ButtonThree = 3

	StickUp    = 101
	StickDown  = 102
	StickLeft  = 103
	StickRight = 104
	StickClick = 105

	KeyPressUp    = 101
	KeyPressDown  = 102
	KeyPressLeft  = 103
	KeyPressRight = 104

	ScreenTouch      = 201
	ScreenSwipeRight = 202
	ScreenSwipeLeft  = 203
	ScreenSwipeUp    = 204
	ScreenSwipeDown  = 205

	EventKindQuit = 0
)

const (
	EventTypeSystem      = 0x1
	EventTypeButton      = 0x2
	EventTypeKey         = 0x3
	EventTypeInput       = 0x4
	EventTypeTouch       = 0x5
	EventTypeApplication = 0x10
	EventTypeUser        = 0x100
)

type Icon struct {
}

type Event struct {
	Type    uint64
	Kind    uint64
	Content interface{}
	Done    func(err error)
}

type EventHandler interface {
	OnEvent(ev *Event) bool
}

type EventListener interface {
	Listen(OnEvent func(ev *Event)) error
}

func NewKeyPressEvent(key int) *Event {
	return &Event{
		Type: EventTypeKey,
		Kind: uint64(key),
	}
}

func NewTouchEvent(touchId int, x, y int) *Event {
	return &Event{
		Type: EventTypeTouch,
		Kind: uint64(touchId),
		Content: &image.Point{ // FIXME: use a dedicatted type?
			X: x,
			Y: y,
		},
	}
}

type SwipeData struct {
	StartX int
	StartY int
	EndX   int
	EndY   int
}

func NewSwipeEvent(direction uint64, touchId int, startX, startY, endX, endY int) *Event {
	//log.Debugf("NewSwipeEvent: %d %d %d.%d %d.%d", direction, touchId, startX, startY, endX, endY)
	return &Event{
		Type: direction,
		Kind: uint64(touchId),
		Content: &SwipeData{
			StartX: startX,
			StartY: startY,
			EndX:   endX,
			EndY:   endY,
		},
	}
}

func NewSystemEvent(kind uint64) *Event {
	return &Event{
		Type: EventTypeSystem,
		Kind: kind,
	}
}

func EventPoint(ev *Event) (image.Point, bool) {
	if ev.Type == EventTypeTouch {
		pt, ok := ev.Content.(*image.Point)
		return *pt, ok
	}

	return image.Point{
		X: -10000, // way off screen
		Y: -10000,
	}, false
}
