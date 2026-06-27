package web

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"github.com/thirdmartini/gogui/pkg/log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Display struct {
	image   *image.RGBA
	address string
}

func defaultView(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `
<html>
<head>
<meta http-equiv="refresh" content="0.25" >
</head>
<img src="/canvas">
</html>
`)
}

func (d *Display) canvasHandler(w http.ResponseWriter, r *http.Request) {
	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, d.image); err != nil {
		log.Println("unable to encode image.")
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
	}
}

func (d *Display) Draw(im *image.RGBA) error {
	draw.Draw(d.image, d.image.Bounds(), im, image.Point{X: 0, Y: 0}, draw.Src)
	return nil
}

func (d *Display) listenAndServe() error {
	r := mux.NewRouter()
	r.HandleFunc("/canvas", d.canvasHandler)
	r.HandleFunc("/", defaultView)
	http.Handle("/", r)
	fmt.Printf("Listening to: http://%s\n", d.address)
	return http.ListenAndServe(d.address, nil)
}

func (d *Display) Size() image.Point {
	return d.image.Bounds().Size()
}

func Open(address string, w, h int) *Display {
	rect := image.Rect(0, 0, w, h)

	wh := &Display{
		image:   image.NewRGBA(rect),
		address: address,
	}

	go func() {
		err := wh.listenAndServe()
		panic(err)
	}()

	return wh
}
