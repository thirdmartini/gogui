package main

import (
	"flag"
	"fmt"
	"image"
	"log"
	"time"

	"github.com/thirdmartini/gogui/pkg/app"
	"github.com/thirdmartini/gogui/pkg/app/views"
	"github.com/thirdmartini/gogui/pkg/app/widgets"
	"github.com/thirdmartini/gogui/pkg/drivers/display/linux/drm"
	"github.com/thirdmartini/gogui/pkg/drivers/input/controller/keyboard"
	"github.com/thirdmartini/gogui/pkg/drivers/input/hid"
	"github.com/thirdmartini/gogui/pkg/ux/themes"

	"github.com/thirdmartini/gogui/pkg/drivers/display"
	"github.com/thirdmartini/gogui/pkg/drivers/display/linux/framebuffer"
	"github.com/thirdmartini/gogui/pkg/drivers/display/vnc"
	"github.com/thirdmartini/gogui/pkg/ux"
	"github.com/thirdmartini/gogui/pkg/ux/canvas"
)

type FrameBufferDevice struct {
	Device   string
	Width    int
	Height   int
	Rotation int
}

func mustInitializeFramebuffers(fbDevices []FrameBufferDevice) ([]display.Display, []ux.EventListener) {
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

func mustInitializeDemo(listenAddress string, width, height int) ([]display.Display, []ux.EventListener) {
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

func main() {
	driverFlag := flag.String("driver", "vnc", "display driver [vnc, framebuffer, drm]")
	touchDeviceFlag := flag.String("touch", "/dev/input/by-id/usb-WaveShare_WaveShare_000000000089-event-if00", "path to touchscreen devicet")
	flag.Parse()

	var displays []display.Display
	var events []ux.EventListener

	switch *driverFlag {
	case "vnc":
		// NOTE: this vnc server is really simple stupid but is compatible with TigerVNC
		// Start a vnc server that will act like the gui display
		// this is nice and useful for testing the ui without a linux FB device
		// note that VNC also provides an event source

		displays, events = mustInitializeDemo(":9000", 1480, 320)
	case "framebuffer":
		// For this demo I'm using a RPI4 with an HDMI touch Display that is 320x1480 but I want to use it
		// as a 1480x320 display
		fbDevs := []FrameBufferDevice{
			{
				Device:   "/dev/fb0",
				Width:    320,
				Height:   1480,
				Rotation: display.Rotation90,
			},
		}
		displays, events = mustInitializeFramebuffers(fbDevs)
		if len(displays) == 0 {
			panic("no framebuffers found  ( try running with --vnc for a demo )")
		}

	case "dri", "drm":
		d, err := drm.NewDisplay("/dev/dri/card1")
		if err != nil {
			panic(err)
		}
		d.WithRotation(display.Rotation90)
		displays = append(displays, d)
	default:
		fmt.Errorf("Unknon driver %s\n", *driverFlag)
	}

	if *touchDeviceFlag != "" {
		touch, err := hid.NewTouchScreen(*touchDeviceFlag)
		touch.SetScaling(1480.0/4000.0, 320.0/4000.0)
		if err != nil {
			log.Printf("Warning: No touch device at %s (Err:%s)\n", *touchDeviceFlag, err)
		} else {
			events = append(events, touch)
		}
	}

	kb, err := keyboard.NewKeyboard()
	if err == nil {
		events = append(events, kb)
	}

	themes.SetTheme("assets/light")
	if err := themes.LoadColors(canvas.NewGGCanvas(image.NewRGBA(image.Rect(0, 0, 1, 1))).ColorPalette()); err != nil {
		panic(err)
	}
	if err := themes.LoadFonts(); err != nil {
		panic(err)
	}

	// create the application
	c := &app.Controller{}

	metric := NewMetrics()

	for idx := range displays {
		p := views.NewPager()
		p.SetBackground(themes.LoadImage("background.png"))
		p.Add(views.NewMainView())
		p.Add(views.NewMeterView().Add(10, 10, &widgets.SpeedGauge{
			Title:       "",
			Width:       800,
			Height:      300,
			LeftMetrics: &metric.EdgeBandwidthIn,
			LeftMax:     1024 * 1024 * 1024,
			LeftLabel:   "DOWNLOAD",

			RightMetrics: &metric.EdgeBandwidthOut,
			RightMax:     1024 * 1024 * 1024,
			RightLabel:   "UPLOAD",

			CenterMetric: metric.EdgePingAvg10s,
			CenterMax:    50.0,
		}))

		c.AddViewPort(views.NewDisplayCanvas(displays[idx]), p)
	}

	app := ux.NewApplication()

	pc, err := NewPingCollector("8.8.8.8")
	go pc.Collect("Edge", func(label string, rtt time.Duration) {
		metric.EdgePingAvg10s.Set(float64(rtt.Microseconds()) / 1000.0)
	})

	go func() {
		downIdx := 0
		upIdx := 10

		for {
			dlNumbers := []float64{100, 100, 150, 200, 200, 100, 100, 500, 0, 100, 200, 300}
			metric.EdgeBandwidthIn.Set(dlNumbers[downIdx] * 1024 * 1024)
			metric.EdgeBandwidthOut.Set(dlNumbers[upIdx] * 1024 * 1024)

			upIdx++
			if upIdx >= len(dlNumbers) {
				upIdx = 0
			}

			downIdx++
			if downIdx >= len(dlNumbers) {
				downIdx = 0
			}

			time.Sleep(time.Second)
		}
	}()

	// start event handler
	for idx := range events {
		go events[idx].Listen(app.PostEvent)
	}

	app.Run(c)
}
