package gdesk

import (
	"testing"
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
	f.Release()
	f.Release()
	c.Stop()
}
