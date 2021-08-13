package gdesk

import (
	"testing"
	"image"
	"image/png"
	"os"
)

func TestCapture(t *testing.T) {
	display := GetPrimaryDisplay()
	var c = NewCapturer()
	c.Start(display)
	for {
		f := c.GetFrame()
		if f != nil {
			break
		}
	}
	f := c.GetFrame()
	w := int(display.Width())
	h := int(display.Height())
	rect := image.Rect(0, 0, w, h)
	rgb := make([]byte, w*h*3)
	
	img := image.NewRGBA(rect)
	i420ToRgb(w, h, f.data, rgb)
	file, _ := os.Create("image.png")
	dst := []byte{}
	stride := len(f.data) / h;
	for y :=0; y < h; y++ {
		for x :=0; x < w; x++ {
			i := stride * y + 3 * x
			dst = append(dst, []byte{rgb[i], rgb[i+1], rgb[i+2], 255}...)
		}
	}
	img.Pix = dst
	png.Encode(file, img)
	f.Release()
	c.Stop()
}
