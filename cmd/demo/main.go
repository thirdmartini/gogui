package main

import (
	"flag"
	"fmt"
	"image"
	"log"
	"time"

	"github.com/thirdmartini/gogui/pkg/app/views"
	"github.com/thirdmartini/gogui/pkg/app/widgets"
	"github.com/thirdmartini/gogui/pkg/drivers/display"
	"github.com/thirdmartini/gogui/pkg/drivers/display/linux/drm"
	"github.com/thirdmartini/gogui/pkg/drivers/display/linux/framebuffer"
	"github.com/thirdmartini/gogui/pkg/drivers/display/vnc"
	"github.com/thirdmartini/gogui/pkg/drivers/input/controller/keyboard"
	"github.com/thirdmartini/gogui/pkg/drivers/input/hid"
	"github.com/thirdmartini/gogui/pkg/ux"
	"github.com/thirdmartini/gogui/pkg/ux/canvas"
	"github.com/thirdmartini/gogui/pkg/ux/canvas/fonts"
	"github.com/thirdmartini/gogui/pkg/ux/themes"
	uxwidget "github.com/thirdmartini/gogui/pkg/ux/widgets"
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

type MainView struct {
	DisplayCanvas *views.DisplayCanvas
	Pager         *views.Pager
}

func (mv *MainView) Next() {
	mv.Pager.Next()
}

func (mv *MainView) OnEvent(event *ux.Event) bool {
	fmt.Printf("Event: %d\n", event.Kind)
	return mv.Pager.OnEvent(event)
}

func (mv *MainView) OnRepaint() {
	mv.Pager.Draw(mv.DisplayCanvas)
	mv.DisplayCanvas.Show()
}

func main() {
	driverFlag := flag.String("driver", "vnc", "display driver [vnc, framebuffer, drm]")
	touchDeviceFlag := flag.String("touch", "/dev/input/by-id/usb-WaveShare_WaveShare_000000000089-event-if00", "path to touchscreen devicet")
	flag.Parse()

	var events []ux.EventListener

	var mainDisplay display.Display

	switch *driverFlag {
	case "vnc":
		// NOTE: this vnc server is really simple stupid but is compatible with TigerVNC
		// Start a vnc server that will act like the gui display
		// this is nice and useful for testing the ui without a linux FB device
		// note that VNC also provides an event source
		var displays []display.Display
		displays, events = mustInitializeDemo(":9000", 1480, 320)
		mainDisplay = displays[0]

	case "framebuffer":
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
		// this forces rotation to 90 degrees, in my case I have a display that is 320x1480, but i want to use
		// it flipped 90degrees and act as a 1480x320
		d.WithRotation(device.Rotation)
		mainDisplay = d

	case "dri", "drm":
		d, err := drm.NewDisplay("/dev/dri/card1")
		if err != nil {
			panic(err)
		}
		d.WithRotation(display.Rotation90)
		defer d.Close()

		mainDisplay = d

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

	metric := NewMetrics()

	p := views.NewPager()
	p.SetBackground(themes.LoadImage("background.png"))

	mv := &MainView{
		DisplayCanvas: views.NewDisplayCanvas(mainDisplay),
		Pager:         p,
	}

	font, err := fonts.Load("assets/light/default.ttf", 30)

	//p.Add(main)
	main := views.NewMainView()
	bt := uxwidget.NewButton(200, 200, 200, 40, ux.AlignLeft, "Click Me", font, themes.ColorTextMuted)
	bt.BackgroundColor = themes.ColorBackground
	bt.BorderColor = themes.ColorBorder
	bt.SetBorder(uxwidget.BorderAll, themes.ColorTextPrimary)
	bt.OnClick = func() bool {
		fmt.Printf("[[Button Clicked]]\n")
		mv.Next()
		return true
	}
	main.AddWidget(bt)
	p.Add(main)

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

	app.Run(mv)
}
