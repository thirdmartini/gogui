package hid

import (
	"encoding/binary"
	"os"
	"syscall"
)

const (
	EV_SYN       = 0x00
	EV_KEY       = 0x01
	EV_REL       = 0x02
	EV_ABS       = 0x03
	EV_MSC       = 0x04
	EV_SW        = 0x05
	EV_LED       = 0x11
	EV_SND       = 0x12
	EV_REP       = 0x14
	EV_FF        = 0x15
	EV_PWR       = 0x16
	EV_FF_STATUS = 0x17

	BTN_TOUCH          = 0x14a // Touch contact
	ABS_X              = 0x00
	ABS_Y              = 0x01
	ABS_Z              = 0x02
	ABS_RX             = 0x03
	ABS_RY             = 0x04
	ABS_RZ             = 0x05
	ABS_MT_SLOT        = 0x2f // Multi-touch slot
	ABS_MT_TOUCH_MAJOR = 0x30
	ABS_MT_POSITION_X  = 0x35
	ABS_MT_POSITION_Y  = 0x36
	ABS_MT_TRACKING_ID = 0x39
)

var inputEventTypeString = map[uint16]string{
	EV_SYN:       "sync",
	EV_REL:       "relative",
	EV_KEY:       "key",
	EV_ABS:       "absolute",
	EV_MSC:       "misc",
	EV_SW:        "switch",
	EV_LED:       "led",
	EV_SND:       "sound",
	EV_REP:       "repeat",
	EV_FF:        "force_feedback",
	EV_PWR:       "power",
	EV_FF_STATUS: "force_feedback_status",
}

var inputEventCodeString = map[uint16]string{
	BTN_TOUCH:          "touch",
	ABS_X:              "x",
	ABS_Y:              "y",
	ABS_Z:              "z",
	ABS_RX:             "rx",
	ABS_RY:             "ry",
	ABS_RZ:             "rz",
	ABS_MT_SLOT:        "slot",
	ABS_MT_TOUCH_MAJOR: "touch_major",
	ABS_MT_POSITION_X:  "position_x",
	ABS_MT_POSITION_Y:  "position_y",
	ABS_MT_TRACKING_ID: "tracking_id",
}

func InputEventTypeString(t uint16) string {
	s, ok := inputEventTypeString[t]
	if !ok {
		panic("unknown input event code")
	}
	return s
}

func InputEventCodeString(t uint16) string {
	s, ok := inputEventCodeString[t]
	if !ok {
		panic("unknown input event code")
	}
	return s
}

type Event struct {
	Time  syscall.Timeval // 16 bytes on 64-bit
	Type  uint16          // event type (EV_SYN, EV_KEY, EV_ABS, etc.)
	Code  uint16          // event code (BTN_TOUCH, ABS_X, ABS_Y, etc.)
	Value int32           // value
}

func (e *Event) TypeString() string {
	return InputEventTypeString(e.Type)
}

func (e *Event) CodeString() string {
	return InputEventCodeString(e.Code)
}

type EventReader struct {
	device *os.File
}

func (r *EventReader) Read() (*Event, error) {
	buf := make([]byte, binary.Size(Event{}))
	_, err := r.device.Read(buf)
	if err != nil {
		return nil, err
	}

	var ev Event
	//binary.LittleEndian.PutUint64(buf[0:8], uint64(ev.Time.Sec)) // Simplified - better to use proper unpacking

	// Let's do it properly:
	//ev.Time.Sec = int64(binary.LittleEndian.Uint64(buf[0:8]))
	//ev.Time.Usec = int32(binary.LittleEndian.Uint64(buf[8:16]))
	ev.Type = binary.LittleEndian.Uint16(buf[16:18])
	ev.Code = binary.LittleEndian.Uint16(buf[18:20])
	ev.Value = int32(binary.LittleEndian.Uint32(buf[20:24]))
	return &ev, nil
}

func (r *EventReader) Close() error {
	return r.device.Close()
}

func NewEventReader(device string) (*EventReader, error) {
	fin, err := os.OpenFile(device, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}

	return &EventReader{
		device: fin,
	}, nil
}
