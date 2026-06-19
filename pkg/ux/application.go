package ux

import (
	"image"
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
}

func (app *Application) PostEvent(ev *Event) {
	app.eventQueue <- ev
}

func (app *Application) WithTheme(themeSource string) *Application {
	themes.SetTheme("assets/light")
	if err := themes.LoadColors(canvas.NewGGCanvas(image.NewRGBA(image.Rect(0, 0, 1, 1))).ColorPalette()); err != nil {
		panic(err)
	}
	if err := themes.LoadFonts(); err != nil {
		panic(err)
	}

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
		case <-time.After(time.Second):
			app.ctrl.OnRepaint()

		// Handle UX or application events
		case ev, ok := <-app.eventQueue:
			if !ok {
				return // channel closed on us, exit the application
			}
			if app.ctrl.OnEvent(ev) {
				app.ctrl.OnRepaint()
			}
		}
	}
}

func (app *Application) Terminate() {
	close(app.eventQueue)
}

func NewApplication() *Application {
	return &Application{
		eventQueue: make(chan *Event),
	}
}
