package gdesk


import (
	"testing"
	"time"
)

func TestCodec(t *testing.T) {
	display := GetPrimaryDisplay()
	var c = NewCapturer()
	c.Start(display)
	var f *Frame
	var cfg = NewEncoderCfg(int(display.Width()), int(display.Height()), VP9)
	var encoder = NewEncoder(cfg)
	writer, _ := NewVpxWriter("aa.webm", cfg)
	ts := time.Now()
	for ii := 0; ii < 100; ii++ {
		f = c.GetFrame()
		if f != nil {
			defer f.Release()
			te := time.Now()
			_ = uint64(te.UnixNano() - ts.UnixNano())
			encodeFrames := encoder.Encode(f, uint64(ii))
			var frame *EncodeFrame
			for frame = encodeFrames.Next(); frame != nil; frame = encodeFrames.Next() {
				writer.WriteFrame(frame)
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
	writer.Close()
}
