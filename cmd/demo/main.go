package main

import (
	"flag"
	"fmt"
	"image"
	"log"

	"github.com/thirdmartini/gogui/pkg/ux"
	"github.com/thirdmartini/gogui/pkg/ux/canvas/fonts"
	"github.com/thirdmartini/gogui/pkg/ux/themes"
	"github.com/thirdmartini/gogui/pkg/ux/widgets"
)

func main() {
	driverFlag := flag.String("driver", "vnc", "display driver [vnc, framebuffer, drm]")
	touchDeviceFlag := flag.String("touch", "/dev/input/by-id/usb-WaveShare_WaveShare_000000000089-event-if00", "path to touchscreen devicet")
	flag.Parse()

	// Initialize the hardware for our display device
	displays, events := mustInitializeHardware(*driverFlag, *touchDeviceFlag)
	if len(displays) == 0 {
		panic("no displays found  ( try running with --vnc for a demo )")
	}
	defer func() {
		for id, d := range displays {
			log.Printf("Closing display %d\n", id)
			d.Close()
		}
	}()

	mainDisplay := displays[0]

	// Load our theme
	// FIXME: make it all happen once

	// make the theme load happen here
	//app, err := ux.NewApplication().WithTheme("assets/light")
	app := ux.NewApplication().WithTheme("assets/light")

	mainWindow := ux.NewWindow(image.Rectangle{
		Max: mainDisplay.Size(),
	})

	mainWindow.SetBackground(themes.LoadImage("background.png"))

	font, _ := fonts.Load("assets/light/default.ttf", 30)

	bt := widgets.NewButton(200, 200, 200, 40, ux.AlignLeft, "Click Me", font, themes.Default.Colors.TextPrimary)
	bt.SetBorder(ux.BorderAll, themes.Default.Colors.TextPrimary)

	bt.OnClick = func() bool {
		fmt.Printf("[[Button Clicked]]\n")
		app.Terminate()
		return true
	}

	mainWindow.AddWidget("button", bt)
	mainWindow.AddWidget("ping", widgets.NewTextLabel(20, 20, ux.AlignLeft, "Click the button to exit", font, themes.Default.Colors.TextPrimary, themes.Default.Colors.Background))

	// start event handler

	dc := ux.NewDisplayController(mainDisplay, mainWindow)
	app.Run(dc, events)
}
