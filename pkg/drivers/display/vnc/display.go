package vnc

import (
	"fmt"
	"image"
	"image/draw"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/thirdmartini/gogui/pkg/drivers/display"
	"github.com/thirdmartini/gogui/pkg/ux"
)

type Display struct {
	address string
	rect    image.Rectangle
	done    chan bool

	image *image.RGBA

	width  int
	height int

	lock   sync.Mutex
	signal *sync.Cond

	rotation int

	pointerX, pointerY int
	buttons            uint8

	events chan interface{}
}

func (r *Display) listenAndServe() error {
	ln, err := net.Listen("tcp", r.address)
	if err != nil {
		return err
	}

	s := NewServer(r.width, r.height)
	go func() {
		err := s.Serve(ln)
		log.Fatalf("rfb server terminated: %v", err)
	}()

	addr := r.address
	if strings.HasPrefix(addr, ":") {
		addr = fmt.Sprintf("localhost%s", addr)
	}

	log.Printf("Listening on: vnc://%s\n", addr)

	for c := range s.Conns {
		r.handleConn(c)
	}
	return nil
}

func (r *Display) Bounds() image.Rectangle {
	switch r.rotation {
	case display.RotationNone, display.Rotation180:
		return r.rect
	case display.Rotation90, display.Rotation270:
		return image.Rect(0, 0, r.rect.Dy(), r.rect.Dx())
	}
	return r.rect
}

func (r *Display) WithRotation(rotation int) *Display {
	r.rotation = rotation
	return r
}

func (r *Display) Close() error {
	// send close signal (FIXME)
	//  right now there is no way to exit here
	return nil
}

func (r *Display) handleConn(c *Conn) { //, rec capture.ImageStream, rect image.Rectangle) {
	rect := image.Rect(0, 0, r.width, r.height)

	im := image.NewRGBA(rect)
	li := &LockableImage{Img: im}

	closec := make(chan bool)
	go func() {
		tick := time.NewTicker(time.Second / 30)
		defer tick.Stop()
		haveNewFrame := false
		for {
			feed := c.Feed
			if !haveNewFrame {
				feed = nil
			}
			_ = feed
			select {
			case feed <- li:
				haveNewFrame = false
			case <-closec:
				return
			case <-tick.C:
				li.Lock()
				draw.Draw(im, im.Bounds(), r.image, r.image.Bounds().Min, draw.Src)
				li.Unlock()
				haveNewFrame = true
			}
		}
	}()

	for e := range c.Event {
		select {
		case r.events <- e:
		default:
			fmt.Println("don't panic")
			log.Printf("ignored event: %#v", e)
		}
	}
	close(closec)
	log.Printf("Client disconnected")
}

func (r *Display) Listen(OnEvent func(ev *ux.Event)) error {
	r.events = make(chan interface{}, 100)
	defer func() {
		r.events = nil
		close(r.events)
	}()

	for e := range r.events {
		switch e.(type) {
		case KeyEvent:
			ve := e.(KeyEvent)
			if ve.DownFlag == 1 {
				if key, ok := remap[ve.Key]; ok {
					ev := &ux.Event{
						Type: ux.EventTypeKey,
						Kind: uint64(key),
					}
					OnEvent(ev)
				} else {
					log.Printf("kb event: %#v", ve)
				}
			}
		case PointerEvent:
			ve := e.(PointerEvent)

			r.pointerX = int(ve.X)
			r.pointerY = int(ve.Y)

			if r.buttons != ve.ButtonMask {
				if ve.ButtonMask == 0 {
					OnEvent(ux.NewTouchEvent(0, r.pointerX, r.pointerY))
				}
				r.buttons = ve.ButtonMask
			}

		default:
		}

	}
	return nil
}

func Rotate90CWFast(dst, src *image.RGBA) {
	b := src.Bounds()
	w := b.Dx()
	h := b.Dy()

	// Stride-aware direct memory copy
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			// Source pixel offset
			srcIdx := (y*src.Stride + x*4)

			// Destination: (y, w-1-x)
			dstX := y
			dstY := w - 1 - x
			dstIdx := (dstY*dst.Stride + dstX*4)

			// Copy 4 bytes (R, G, B, A)
			copy(dst.Pix[dstIdx:dstIdx+4], src.Pix[srcIdx:srcIdx+4])
		}
	}
}

func (r *Display) Draw(im *image.RGBA) error {
	//fmt.Printf("src: %dx%d    dst: %dx%d\n", im.Bounds().Dx(), im.Bounds().Dy(), r.image.Bounds().Dx(), r.image.Bounds().Dy())

	switch r.rotation {
	case display.RotationNone, display.Rotation180:
		draw.Draw(r.image, r.image.Bounds(), im, image.Point{X: 0, Y: 0}, draw.Src)
	case display.Rotation90, display.Rotation270:
		Rotate90CWFast(r.image, im)
	}

	r.signal.Broadcast()
	return nil
}

func (r *Display) Size() image.Point {
	switch r.rotation {
	case display.RotationNone, display.Rotation180:
		return image.Point{
			X: r.width,
			Y: r.height,
		}
	case display.Rotation90, display.Rotation270:
		return image.Point{
			X: r.height,
			Y: r.width,
		}
	}

	return image.Point{
		X: r.width,
		Y: r.height,
	}
}

func Open(address string, width, height int) (*Display, error) {
	rect := image.Rect(0, 0, width, height)
	im := image.NewRGBA(rect)

	r := &Display{
		address: address,
		width:   width,
		height:  height,
		image:   im,
	}
	r.signal = sync.NewCond(&r.lock)

	go r.listenAndServe()

	return r, nil
}
