package hid

import (
	"fmt"
	"log"

	"github.com/thirdmartini/gogui/pkg/ux"
)

const (
	TOUCH_DOWN = 1
	TOUCH_UP   = 2
)

type TouchScreen struct {
	scaleX float64
	scaleY float64
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

		//log.Printf("%s:%s:: %+v\n", ev.TypeString(), ev.CodeString(), ev)

		switch ev.Type {
		case EV_KEY:
			switch ev.Code {
			case BTN_TOUCH:
				e <- TouchEvent{Type: ev.Value, X: CurrentX, Y: CurrentY}
			default:
				//log.Printf("EV_KEY:unknown key event %d", ev.Code)
			}

		case EV_ABS:
			switch ev.Code {
			case ABS_X, ABS_MT_POSITION_X:
				CurrentX = ev.Value

			case ABS_Y, ABS_MT_POSITION_Y:
				CurrentY = ev.Value

			default:
				//log.Printf("EV_ABS:unknown key event", ev.Code)
			}
		case EV_SYN:
			//log.Printf("EV_SYNC:unknown key event", ev.Code)

		default:
			//log.Printf("unknown event type", ev.Type)
		}
	}
}

func (t *TouchScreen) SetScaling(scaleX, scaleY float64) {
	t.scaleX = scaleX
	t.scaleY = scaleY
}

func (t *TouchScreen) Listen(OnEvent func(ev *ux.Event)) error {
	c := make(chan TouchEvent, 10)
	go func() {
		_ = t.poll(c)
		close(c)
	}()

	for te := range c {
		x := int(float64(te.X) * t.scaleX)
		y := int(float64(te.Y) * t.scaleY)

		log.Printf("E: %+v [%v/%v]", te, t.scaleX, t.scaleY)

		// normalize these
		OnEvent(ux.NewTouchEvent(0, x, y))
	}
	return fmt.Errorf("exited")
}

func (t *TouchScreen) Close() error {
	return t.EventReader.Close()
}

func NewTouchScreen(device string) (*TouchScreen, error) {
	reader, err := NewEventReader(device)
	if err != nil {
		return nil, err
	}
	return &TouchScreen{
		EventReader: reader,
		scaleY:      1.0,
		scaleX:      1.0,
	}, nil
}
