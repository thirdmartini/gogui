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
	"github.com/thirdmartini/gogui/pkg/drivers/display"
	"github.com/thirdmartini/gogui/pkg/drivers/display/linux/drm"
	"github.com/thirdmartini/gogui/pkg/drivers/input/controller/keyboard"
	"github.com/thirdmartini/gogui/pkg/drivers/input/hid"
	"github.com/thirdmartini/gogui/pkg/ux"
	"github.com/thirdmartini/gogui/pkg/ux/canvas"
	"github.com/thirdmartini/gogui/pkg/ux/canvas/fonts"
	"github.com/thirdmartini/gogui/pkg/ux/themes"
	uxwidget "github.com/thirdmartini/gogui/pkg/ux/widgets"
)

func mainMultiDisplay() {
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
		if err != nil {
			log.Printf("Warning: No touch device at %s (Err:%s)\n", *touchDeviceFlag, err)
		} else {
			touch.SetScaling(1480.0/4000.0, 320.0/4000.0)
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

		//p.Add(main)
		main := views.NewMainView()
		bt := uxwidget.NewButton(200, 200, 100, 50, ux.AlignLeft, "Click Me", fonts.Font16Internal, themes.ColorTextMuted)
		bt.BackgroundColor = themes.ColorMenuBackground
		bt.BorderColor = themes.ColorBorder
		bt.OnClick = func() bool {
			fmt.Printf("[[Button Clicked]]\n")
			c.Next()
			return true
		}

		main.AddWidget(bt)

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

		// viewports handle multiple displays and switching between which display receives input
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
