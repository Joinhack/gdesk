package gdesk

import (
	"os"
	"encoding/binary"
	"bytes"
	"io"
)

type VpxWriter struct {
	file *os.File
	cfg *EncoderCfg
	count int
}

const VP8_FOURCC = 0x30385056
const VP9_FOURCC = 0x30395056

func NewVpxWriter(fileName string, cfg *EncoderCfg) (*VpxWriter, error) {
	var writer VpxWriter
	var err error
	if writer.file, err = os.Create(fileName); err != nil {
		return nil, err
	}
	writer.cfg = cfg
	writer.writeHeader()
	return &writer, nil
}

func (w *VpxWriter) writeHeader() {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, []byte("DKIF"))
	binary.Write(buf, binary.LittleEndian, uint16(0))
	binary.Write(buf, binary.LittleEndian, uint16(32))
	
	var fourcc int32 = VP8_FOURCC
	if w.cfg.Codec == VP9 {
		fourcc = VP9_FOURCC
	}
	var cfg = w.cfg
	binary.Write(buf, binary.LittleEndian, fourcc)
	binary.Write(buf, binary.LittleEndian, uint16(cfg.Width))
	binary.Write(buf, binary.LittleEndian, uint16(cfg.Height))
	binary.Write(buf, binary.LittleEndian, uint32(cfg.Denominator))
	binary.Write(buf, binary.LittleEndian, uint32(cfg.Numerator))
	binary.Write(buf, binary.LittleEndian, uint32(w.count))
	binary.Write(buf, binary.LittleEndian, uint32(0))
	w.file.Write(buf.Bytes())
}

func (w *VpxWriter) writeFrameHeader(len int, pts uint64) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, uint32(len))
	binary.Write(buf, binary.LittleEndian, uint32(pts&0xFFFFFFFF))
	binary.Write(buf, binary.LittleEndian, uint32(pts>>32))
	w.file.Write(buf.Bytes())
}

func (w *VpxWriter) WriteFrame(f *EncodeFrame) {
	w.writeFrameHeader(len(f.Data), f.Pts)
	w.file.Write(f.Data)
	w.count++
}

func (w *VpxWriter) Close() {
	if w.file != nil {
		if _, err := w.file.Seek(0, io.SeekStart); err != nil {
			panic(err)
		}
		w.writeHeader()
		w.file.Close()
		w.file = nil
	}
}