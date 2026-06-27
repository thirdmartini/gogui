package ux

import (
	"fmt"
	"image"
	"github.com/thirdmartini/gogui/pkg/log"
	"time"

	"github.com/thirdmartini/gogui/pkg/ux/canvas"
	"github.com/thirdmartini/gogui/pkg/ux/themes"
)

type ViewController interface {
	OnEvent(ev *Event) bool
	OnRepaint()
}

type ApplicationController interface {
	OnEvent(ev *Event) bool
	OnRepaint()
}

type Application struct {
	eventQueue chan *Event
	ctrl       ApplicationController
	interval   time.Duration
}

func (app *Application) PostEvent(ev *Event) {
	app.eventQueue <- ev
}

func (app *Application) WithTheme(themeSource string) *Application {
	palette := canvas.NewGGCanvas(image.NewRGBA(image.Rect(0, 0, 1, 1))).ColorPalette()

	theme, err := themes.Load(themeSource, palette)
	if err != nil {
		panic(fmt.Sprintf("failed to load theme %s", themeSource))
	}
	themes.SetTheme(theme)
	return app
}

func (app *Application) WithRefreshRate(hz uint) *Application {
	if hz == 0 {
		return app
	}

	if hz > 1000 {
		hz = 1000
	}

	app.interval = time.Duration(1000/hz) * time.Millisecond
	return app
}

func (app *Application) Run(ctrl ViewController, eventSources []EventListener) {

	for idx := range eventSources {
		go func() {
			err := eventSources[idx].Listen(app.PostEvent)
			if err != nil {
				// FIXME post a clean exit event to queue
				panic(err)
			}
		}()
	}

	app.ctrl = ctrl
	for {
		select {
		// Screen update interval
		case <-time.After(app.interval):
			app.ctrl.OnRepaint()

		// Handle UX or application events
		case ev, ok := <-app.eventQueue:
			if !ok {
				return // channel closed on us, exit the application
			}

			switch ev.Type {
			case EventTypeSystem:
				switch ev.Kind {
				case EventKindQuit:
					log.Printf("Quitting application")
					return
				}
			}

			if app.ctrl.OnEvent(ev) {
				app.ctrl.OnRepaint()
			}
		}
	}
}

func (app *Application) Terminate() {
	app.PostEvent(NewSystemEvent(EventKindQuit))
}

func NewApplication() *Application {
	return &Application{
		interval:   time.Second,
		eventQueue: make(chan *Event, 128),
	}
}
