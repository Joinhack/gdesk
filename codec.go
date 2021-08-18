package gdesk

/*
#cgo pkg-config: vpx
#include "vpx_wrapper.h"
*/
import "C"
import (
	"reflect"
	"runtime"
	"sync/atomic"
	"unsafe"
)

const (
	VP8 = iota
	VP9
)

type Encoder struct {
	cfg    C.vpx_codec_enc_cfg_t
	codec  C.vpx_codec_ctx_t
	inited int32
}

type EncodeFrames struct {
	codec *C.vpx_codec_ctx_t
	iter  C.vpx_codec_iter_t
	pts   uint64
}

type EncoderCfg struct {
	Numerator   int
	Denominator int
	Bitrate     int
	Width       int
	Height      int
	Codec       int
}

func NewEncoderCfg(w, h, c int) *EncoderCfg {
	return &EncoderCfg{
		Numerator: 1,
		Denominator: 1000,
		Bitrate: 1000,
		Width: w,
		Height: h,
		Codec: c,
	}
}

func NewEncoder(cfg *EncoderCfg) *Encoder {
	var iface *C.vpx_codec_iface_t
	encoder := &Encoder{}

	if cfg.Numerator == 0 {
		cfg.Numerator = 1
	}

	if cfg.Denominator == 0 {
		cfg.Denominator = 1000
	}

	if cfg.Codec == VP8 {
		iface = C.vpx_codec_vp8_cx()
	} else {
		iface = C.vpx_codec_vp9_cx()
	}
	if cfg.Bitrate == 0 {
		cfg.Bitrate = 5000
	}
	if int(C.vpx_codec_enc_config_default(iface, &encoder.cfg, C.uint(0))) != 0 {
		panic("error vpx_codec_enc_config_default")
	}
	encoder.cfg.g_w = C.uint(cfg.Width)
	encoder.cfg.g_h = C.uint(cfg.Height)
	encoder.cfg.g_timebase.num = C.int(cfg.Numerator)
	encoder.cfg.g_timebase.den = C.int(cfg.Denominator)
	encoder.cfg.rc_target_bitrate = C.uint(cfg.Bitrate)
	encoder.cfg.rc_undershoot_pct = 95
	encoder.cfg.rc_dropframe_thresh = 25
	encoder.cfg.g_error_resilient = C.VPX_ERROR_RESILIENT_DEFAULT
	encoder.cfg.rc_end_usage = C.VPX_CBR
	encoder.cfg.kf_mode = C.VPX_KF_DISABLED

	if int(C.vpx_codec_enc_init_ver(&encoder.codec, iface, &encoder.cfg, 0, C.VPX_ENCODER_ABI_VERSION)) != 0 {
		panic("error vpx_codec_enc_init_ver")
	}
	encoder.inited = 1
	runtime.SetFinalizer(encoder, (*Encoder).Release)
	return encoder
}

func (e *Encoder) Release() {
	if atomic.LoadInt32(&e.inited) == 0 {
		return
	}
	for {
		if atomic.CompareAndSwapInt32(&e.inited, e.inited, 0) {
			break
		}
	}
	C.vpx_codec_destroy(&e.codec)
}

func (e *Encoder) Encode(f *Frame, pts uint64) *EncodeFrames {
	d := unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&f.data)).Data)
	if int(C.encode(&e.codec, d, C.int(e.cfg.g_w), C.int(e.cfg.g_h), C.longlong(pts))) != 0 {
		panic("vpx_codec_encode error")
	}
	return &EncodeFrames{
		codec: &e.codec,
		iter:  nil,
		pts:   pts,
	}
}

type EncodeFrame struct {
	Data []byte
	Key  int
	Pts  uint64
}

func (e *EncodeFrames) Next() *EncodeFrame {
	frame := &EncodeFrame{}
	var dataH reflect.SliceHeader
	key := (*C.int)(unsafe.Pointer(&frame.Key))
	data := (**C.char)(unsafe.Pointer(&dataH.Data))
	l := (*C.int)(unsafe.Pointer(&dataH.Len))
	pts := (*C.int)(unsafe.Pointer(&frame.Pts))
	rs := int(C.frame_data(e.codec, &e.iter, key, data, l, pts))
	dataH.Cap = dataH.Len
	if  rs != 0 {
		return nil
	}
	frame.Data = *(*[]byte)(unsafe.Pointer(&dataH))
	return frame

}


