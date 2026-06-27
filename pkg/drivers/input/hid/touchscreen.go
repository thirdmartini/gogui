package hid

import (
	"fmt"

	"github.com/thirdmartini/gogui/pkg/log"

	"github.com/thirdmartini/gogui/pkg/ux"
)

const (
	TOUCH_DOWN = 1
	TOUCH_UP   = 2
)

type TouchScreen struct {
	swipeThresholdX int
	swipeThresholdY int
	scaleX          float64
	scaleY          float64
	*EventReader
}

type TouchEvent struct {
	Type int32
	X    int32
	Y    int32
}

func (t *TouchEvent) String() string {
	switch t.Type {
	case TOUCH_DOWN:
		return "DOWN"
	case TOUCH_UP:
		return "UP"
	default:
		return "UNKNOWN"
	}
}

func (t *TouchScreen) poll(e chan TouchEvent) error {
	CurrentX := int32(0)
	CurrentY := int32(0)

	for {
		ev, err := t.EventReader.Read()
		if err != nil {
			return err
		}

		//log.Debugf("%s:%s:: %+v\n", ev.TypeString(), ev.CodeString(), ev)

		switch ev.Type {
		case EV_KEY:
			switch ev.Code {
			case BTN_TOUCH:
				e <- TouchEvent{Type: ev.Value, X: CurrentX, Y: CurrentY}
			default:
				log.Warnf("EV_KEY:unknown key event %d", ev.Code)
			}

		case EV_ABS:
			switch ev.Code {
			case ABS_X, ABS_MT_POSITION_X:
				CurrentX = ev.Value

			case ABS_Y, ABS_MT_POSITION_Y:
				CurrentY = ev.Value

			default:
				log.Debugf("EV_ABS:unknown key event", ev.Code)
			}
		case EV_SYN:
			log.Debugf("EV_SYNC:unknown key event", ev.Code)

		default:
			log.Debugf("unknown event type", ev.Type)
		}
	}
}

func (t *TouchScreen) SetScaling(scaleX, scaleY float64) {
	t.scaleX = scaleX
	t.scaleY = scaleY
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (t *TouchScreen) Listen(OnEvent func(ev *ux.Event)) error {
	c := make(chan TouchEvent, 10)
	go func() {
		_ = t.poll(c)
		close(c)
	}()

	downX := 0
	downY := 0

	for te := range c {
		x := int(float64(te.X) * t.scaleX)
		y := int(float64(te.Y) * t.scaleY)

		// no actionon down
		if te.Type == 1 {
			downX = x
			downY = y
			continue
		}

		// figure out if this was a swipe event
		swipeX := downX - x
		swipeY := downY - y

		if abs(swipeX) < t.swipeThresholdX {
			swipeX = 0
		}

		if abs(swipeY) < t.swipeThresholdY {
			swipeY = 0
		}

		if abs(swipeX) > abs(swipeY) {
			if swipeX > 0 {
				OnEvent(ux.NewSwipeEvent(ux.ScreenSwipeLeft, 0, downX, downY, x, y))
			} else {
				OnEvent(ux.NewSwipeEvent(ux.ScreenSwipeRight, 0, downX, downY, x, y))
			}
		} else if abs(swipeY) > abs(swipeX) {
			if swipeY > 0 {
				OnEvent(ux.NewSwipeEvent(ux.ScreenSwipeUp, 0, downX, downY, x, y))
			} else {
				OnEvent(ux.NewSwipeEvent(ux.ScreenSwipeDown, 0, downX, downY, x, y))
			}

		} else {
			OnEvent(ux.NewTouchEvent(0, x, y))
		}

		// normalize these

	}
	return fmt.Errorf("exited")
}

func (t *TouchScreen) Close() error {
	return t.EventReader.Close()
}

func (t *TouchScreen) SetSwipeThreshold(thresholdX, thresholdY int) {
	t.swipeThresholdX = thresholdX
	t.swipeThresholdY = thresholdY
}

func NewTouchScreen(device string) (*TouchScreen, error) {
	reader, err := NewEventReader(device)
	if err != nil {
		return nil, err
	}
	return &TouchScreen{
		EventReader:     reader,
		scaleY:          1.0,
		scaleX:          1.0,
		swipeThresholdX: 50, // 50 px
		swipeThresholdY: 50,
	}, nil
}
