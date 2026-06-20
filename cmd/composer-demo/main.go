package main

import (
	"flag"
	"log"

	"github.com/thirdmartini/gogui/cmd"
	"github.com/thirdmartini/gogui/pkg/composer"
	"github.com/thirdmartini/gogui/pkg/ux"
	"github.com/thirdmartini/gogui/pkg/ux/widgets"
)

func main() {
	driverFlag := flag.String("driver", "vnc", "display driver [vnc, framebuffer, drm]")
	touchDeviceFlag := flag.String("touch", "/dev/input/by-id/usb-WaveShare_WaveShare_000000000089-event-if00", "path to touchscreen devicet")
	flag.Parse()

	// Initialize the hardware for our display device
	displays, events := cmd.MustInitializeHardware(*driverFlag, *touchDeviceFlag)
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

	app, compose, err := composer.From("cmd/composer-demo/composer.json")
	if err != nil {
		panic(err)
	}

	compose.Get("button.quit").(*widgets.Button).OnClick = func() bool {
		app.Terminate()
		return true
	}

	dc := ux.NewDisplayController(mainDisplay, compose.Root(0))
	app.Run(dc, events)
}
