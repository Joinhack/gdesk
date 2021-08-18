package gdesk

/*
#cgo LDFLAGS: -lyuv
#include "vpx_wrapper.h"
*/
import "C"

import (
	"reflect"
	"unsafe"
)

func i420ToRgb(w, h int, src []byte) (dst []byte) {
	dst = make([]byte, w*h*3)
	s := unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&src)).Data)
	d := unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&dst)).Data)
	C.i420_to_rgb(C.int(w), C.int(h), s, d)
	return
}

func nv12_to_i420(plan0, plan1 []byte, w, h int) (dst []byte) {
	dst = make([]byte, h*w*12/8)
	p0 := (*reflect.SliceHeader)(unsafe.Pointer(&plan0))
	p1 := (*reflect.SliceHeader)(unsafe.Pointer(&plan1))
	d := unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&dst)).Data)
	C.nv12_to_i420(unsafe.Pointer(p0.Data), C.int(p0.Len), unsafe.Pointer(p1.Data), C.int(p1.Len), C.int(w), C.int(h), d)
	return
}
