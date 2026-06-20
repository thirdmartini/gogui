package cmd

import (
	"fmt"
	"log"

	"github.com/thirdmartini/gogui/pkg/drivers/display"
	"github.com/thirdmartini/gogui/pkg/drivers/display/linux/drm"
	"github.com/thirdmartini/gogui/pkg/drivers/display/linux/framebuffer"
	"github.com/thirdmartini/gogui/pkg/drivers/display/vnc"
	"github.com/thirdmartini/gogui/pkg/drivers/input/controller/keyboard"
	"github.com/thirdmartini/gogui/pkg/drivers/input/hid"
	"github.com/thirdmartini/gogui/pkg/ux"
)

type FrameBufferDevice struct {
	Device   string
	Width    int
	Height   int
	Rotation int
}

func MustInitializeFramebuffers(fbDevices []FrameBufferDevice) ([]display.Display, []ux.EventListener) {
	var displays []display.Display
	var events []ux.EventListener

	for _, device := range fbDevices {
		d, err := framebuffer.Open(device.Device, device.Width, device.Height)
		if err != nil {
			panic(err)
		}
		// this forces rotation to 90 degrees, in my case I have a display that is 320x1480, but i want to use
		// it flipped 90degrees and act as a 1480x320
		d.WithRotation(device.Rotation)
		displays = append(displays, d)
	}

	fmt.Printf("FB Acquired\n")
	return displays, events
}

func MustInitializeDemo(listenAddress string, width, height int) ([]display.Display, []ux.EventListener) {
	var displays []display.Display
	var events []ux.EventListener

	vnc, err := vnc.Open(listenAddress, width, height)
	if err != nil {
		panic(err)
	}

	//
	//vnc.WithRotation(display.Rotation90)

	events = append(events, vnc)
	displays = append(displays, vnc)

	return displays, events
}

func MustInitializeHardware(displayDriver, touchDriver string) (displays []display.Display, events []ux.EventListener) {
	switch displayDriver {
	case "vnc":
		// NOTE: this vnc server is really simple stupid but is compatible with TigerVNC
		// Start a vnc server that will act like the gui display
		// this is nice and useful for testing the ui without a linux FB device
		// note that VNC also provides an event source

		displays, events = MustInitializeDemo(":9000", 1480, 320)

	case "framebuffer":
		// For this demo I'm using a RPI4 with an HDMI touch Display that is 320x1480
		device := FrameBufferDevice{
			Device:   "/dev/fb0",
			Width:    320,
			Height:   1480,
			Rotation: display.Rotation90,
		}

		d, err := framebuffer.Open(device.Device, device.Width, device.Height)
		if err != nil {
			panic(err)
		}
		// this forces rotation to 90 degrees, in my case I have a display that is 320x1480, but I want to use
		// it flipped 90degrees and act as a 1480x320
		d.WithRotation(device.Rotation)
		displays = append(displays, d)

	case "dri", "drm":
		d, err := drm.NewDisplay("/dev/dri/card1")
		if err != nil {
			panic(err)
		}
		d.WithRotation(display.Rotation90)

		displays = append(displays, d)
	default:
		fmt.Errorf("Unknonw driver %s\n", displayDriver)
	}

	if touchDriver != "" {
		touch, err := hid.NewTouchScreen(touchDriver)
		if err != nil {
			log.Printf("Warning: No touch device at %s (Err:%s)\n", touchDriver, err)
		} else {
			touch.SetScaling(1480.0/4000.0, 320.0/4000.0)
			touch.SetSwipeThreshold(100, 100)
			events = append(events, touch)
		}
	}

	kb, err := keyboard.NewKeyboard()
	if err == nil {
		events = append(events, kb)
	}

	return displays, events
}
