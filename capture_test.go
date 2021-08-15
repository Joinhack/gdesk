package gdesk

import (
	"image"
	"image/png"
	"os"
	"testing"
	"time"
)

func TestCapture(t *testing.T) {
	display := GetPrimaryDisplay()
	var c = NewCapturer()
	c.Start(display)
	var f *Frame
	for {
	for {
		f = c.GetFrame()
		if f != nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	w := int(display.Width())
	h := int(display.Height())
	rect := image.Rect(0, 0, w, h)

	img := image.NewRGBA(rect)
	rgb := i420ToRgb(w, h, f.data)
	file, _ := os.Create("image.png")
	dst := make([]byte, w*h*4)
	stride := len(rgb) / h
	idx := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			i := stride*y + 3*x
			dst[idx] = rgb[i]
			dst[idx+1] = rgb[i+1]
			dst[idx+2] = rgb[i+2]
			dst[idx+3] = 255
			idx += 4
		}
	}
	img.Pix = dst
	png.Encode(file, img)
	f.Release()
}
	time.Sleep(20 * time.Millisecond)
	c.Stop()
}
