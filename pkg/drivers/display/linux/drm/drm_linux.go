//go:build linux

package drm

// #cgo CFLAGS: -I/usr/include/libdrm
// #cgo LDFLAGS:  -ldrm
/*
#include <stdlib.h>
#include <stdio.h>
#include <unistd.h>
#include <fcntl.h>
#include <errno.h>
#include <string.h>
#include <sys/poll.h>
#include <sys/mman.h>
#include <xf86drm.h>
#include <xf86drmMode.h>


typedef struct {
    int fd;
	uint32_t connectorId;
    drmModeConnector* conn;
    drmModeRes* res;
	drmModeCrtc* crtc;
	drmModeModeInfo mode;
}display_t;


display_t* drmAllocateDisplay() {
	display_t *display = (display_t*) malloc(sizeof(display_t));
	display->fd = -1;
	display->conn = NULL;
	display->res = NULL;
	display->crtc = NULL;
	return display;
}

void drmReleaseDisplay(display_t *display) {
	if ( display == NULL ) return;

	if (display->res) {
        drmModeFreeResources(display->res);
		display->res = NULL;

	}
	if (display->conn) {
        drmModeFreeConnector(display->conn);
		display->conn = NULL;
	}


	if (display->fd >= 0) {
		close(display->fd);
	}
	display->fd = -1;

	free(display);
}

display_t *drmAcquireDisplay(const char * deviceName) {
	display_t *display = drmAllocateDisplay();

    // 1. Open DRM device
    display->fd = open(deviceName, O_RDWR | O_CLOEXEC);
    if (display->fd < 0) {
		drmReleaseDisplay(display);
        printf("Failed to open DRM device");
        return NULL;
    }

    uint64_t has_dumb;
    if (drmGetCap(display->fd, DRM_CAP_DUMB_BUFFER, &has_dumb) < 0 ) {
        fprintf(stderr, "No dumb buffer support\n");
        drmReleaseDisplay(display);
        return NULL;
    }

	drmSetMaster(display->fd);

    // 2. Get resources and find a connected connector
    display->res = drmModeGetResources(display->fd);
    if (!display->res) {
        drmReleaseDisplay(display);
        return NULL;
    }

    for (int i = 0; i < display->res->count_connectors; ++i) {
        drmModeConnector* conn = drmModeGetConnector(display->fd, display->res->connectors[i]);
        if (conn) {
			printf("Found connected connector: %u (%s) [%d]\n", conn->connector_id,
				   conn->connector_type == DRM_MODE_CONNECTOR_HDMIA ? "HDMI" : "Other", conn->connection);

            if ((conn->connection == DRM_MODE_CONNECTED) && (conn->count_modes > 0 )) {
                display->conn = conn;
				display->mode = conn->modes[0];
                printf("Found connected connector: %u (%s) modes:%d\n", conn->connector_id,
                       conn->connector_type == DRM_MODE_CONNECTOR_HDMIA ? "HDMI" : "Other", conn->count_modes);
				break;
            }
            drmModeFreeConnector(conn);
        }
    }

    if (display->conn == NULL) {
        printf("No connected connector found\n");
        drmReleaseDisplay(display);
        return NULL;
    }

    drmModeEncoder *enc = NULL;
    if (display->conn->encoder_id) {
        enc = drmModeGetEncoder(display->fd, display->conn->encoder_id);
	}

    uint32_t crtc_id = 0;
    if (enc) {
        crtc_id = enc->crtc_id;
        drmModeFreeEncoder(enc);
    }
    if (!crtc_id) {
        // Pick first available CRTC (simple path)
        if (display->res->count_crtcs > 0) {
            crtc_id = display->res->crtcs[0];
        }
    }

    display->crtc = drmModeGetCrtc(display->fd, crtc_id);
    if (!display->crtc) {
        printf("No Crtc found\n");
		drmReleaseDisplay(display);
        return NULL;
    }

    printf("DRM fd=%d Mode %ux%u @ %uHz\n", display->fd, display->mode.vdisplay, display->mode.hdisplay, display->mode.vrefresh);
    return display;
}


typedef struct {
	int fb_id;
	uint8_t *mem;
	uint64_t length;
}framebuffer_t;


int drmInitFramebuffer(display_t *display, framebuffer_t *fb ) {
    struct drm_mode_create_dumb creq = {
        .width = display->mode.hdisplay,
        .height = display->mode.vdisplay,
        .bpp = 32,
    };
    drmIoctl(display->fd, DRM_IOCTL_MODE_CREATE_DUMB, &creq);

    // Add framebuffer
    drmModeAddFB(display->fd, creq.width, creq.height, 24, 32, creq.pitch, creq.handle, &fb->fb_id);

    // Map it
    struct drm_mode_map_dumb mreq = { .handle = creq.handle };
    drmIoctl(display->fd, DRM_IOCTL_MODE_MAP_DUMB, &mreq);

	fb->length = creq.size;
    fb->mem = mmap(0, fb->length, PROT_READ | PROT_WRITE, MAP_SHARED, display->fd, mreq.offset);

	return 0;
}

int drmModeSet(display_t *display, int fb_id) {
	int ret = drmModeSetCrtc(display->fd, display->crtc->crtc_id, fb_id,
                                 0, 0, &display->conn->connector_id, 1, &display->mode);
	if (ret) {
		fprintf(stderr, "drmModeSetCrtc initial failed: %s\n", strerror(-ret));
	}
	return ret;
}

int drmPageFlip(display_t *display, int fb_id) {
	int ret = drmModePageFlip(display->fd, display->crtc->crtc_id, fb_id, DRM_MODE_PAGE_FLIP_EVENT, NULL);
	if (ret) {
    	drmModeSetCrtc(display->fd, display->crtc->crtc_id, fb_id,
                               0, 0, &display->conn->connector_id, 1, &display->mode);
	} else {
		struct pollfd pfd = { .fd = display->fd, .events = POLLIN };
		poll(&pfd, 1, 16);
		drmEventContext ev = { .version = DRM_EVENT_CONTEXT_VERSION };
		drmHandleEvent(display->fd, &ev);
	}
	return 0;
}
*/
import "C"

import (
	"fmt"
	"image"
	"unsafe"

	"github.com/thirdmartini/gogui/pkg/log"

	"github.com/thirdmartini/gogui/pkg/drivers/display"
)

type FrameBuffer C.framebuffer_t

type Display struct {
	display *C.display_t

	Width  int
	Height int

	frameBufferCount int // numbe rof frame buffers to use for mode flipping (default don't flip)
	current          int // current framebuffer
	fbs              []FrameBuffer

	//fd       int
	rotation int
}

func (d *Display) createFrameBuffers() error {
	d.fbs = make([]FrameBuffer, d.frameBufferCount)
	for idx := range d.fbs {
		rc := C.drmInitFramebuffer(d.display, (*C.framebuffer_t)(&d.fbs[idx]))
		if rc != 0 {
			panic("failed to initialize frame buffers")
		}

		log.Debugf("FB:%d -> %d %v %v\n", idx, d.fbs[idx].fb_id, d.fbs[idx].mem, d.fbs[idx].length)
	}
	return nil
}

func (d *Display) acquireDisplay(deviceNode string) error {
	var err error

	s := C.CString(deviceNode)
	display := C.drmAcquireDisplay(s)
	C.free(unsafe.Pointer(s))
	if display == nil {
		return fmt.Errorf("unable to acquire dri/drm display")
	}
	d.display = display
	d.Width = int(display.mode.hdisplay)
	d.Height = int(display.mode.vdisplay)

	err = d.createFrameBuffers()
	if err != nil {
		C.drmReleaseDisplay(display)
	}

	buffer := unsafe.Slice((*byte)(unsafe.Pointer(d.fbs[d.current].mem)), d.fbs[d.current].length)
	for i := 0; i < len(buffer); i++ {
		buffer[i] = 128
	}

	rc := C.drmModeSet(d.display, d.fbs[d.current].fb_id)
	if rc < 0 {
		return fmt.Errorf("failed to mode set display")
	}

	return nil
}

func (d *Display) Close() error {
	if d.display != nil {
		C.drmReleaseDisplay(d.display)
		d.display = nil
	}
	return nil
}

func (d *Display) onDrawReady(drawFunc func(buffer []byte)) {
	buffer := unsafe.Slice((*byte)(unsafe.Pointer(d.fbs[d.current].mem)), d.fbs[d.current].length)
	drawFunc(buffer)
	C.drmPageFlip(d.display, d.fbs[d.current].fb_id)
	d.current = (d.current + 1) % int(d.frameBufferCount)
}

func (d *Display) Draw(im *image.RGBA) error {
	d.onDrawReady(func(buffer []byte) {
		rgb := im.Pix

		switch d.rotation {
		case display.RotationNone, display.Rotation180:
			for pos := 0; pos < len(buffer); pos += 4 {
				buffer[pos] = rgb[pos+2]
				buffer[pos+1] = rgb[pos+1]
				buffer[pos+2] = rgb[pos]
				buffer[pos+3] = rgb[pos+3]
			}
		case display.Rotation90, display.Rotation270:
			// FIXME only do 90
			b := im.Bounds()
			w := b.Dx()
			h := b.Dy()

			stride := d.Width * 4 // we assume 32bit color, might want to muck with that
			for y := 0; y < h; y++ {
				for x := 0; x < w; x++ {
					srcIdx := (y*im.Stride + x*4)

					// Destination: (y, w-1-x)
					dstX := y
					dstY := w - 1 - x
					dstIdx := (dstY*stride + dstX*4)

					buffer[dstIdx] = rgb[srcIdx+2]
					buffer[dstIdx+1] = rgb[srcIdx+1]
					buffer[dstIdx+2] = rgb[srcIdx]
					buffer[dstIdx+3] = rgb[srcIdx+3]
				}
			}
		}
	})

	return nil
}

func (d *Display) WithRotation(rotation int) *Display {
	d.rotation = rotation
	return d
}

func (d *Display) Size() image.Point {
	switch d.rotation {
	case display.Rotation90, display.Rotation270:
		return image.Point{
			X: d.Height,
			Y: d.Width,
		}
	}

	return image.Point{
		X: d.Width,
		Y: d.Height,
	}
}

func NewDisplay(device string) (*Display, error) {
	d := &Display{
		frameBufferCount: 1,
	}

	err := d.acquireDisplay(device)
	if err != nil {
		return nil, err
	}
	return d, nil
}
